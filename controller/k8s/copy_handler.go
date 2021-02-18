package k8s

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"io"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubecp/config"
	"kubecp/controller"
	"kubecp/logs"
	"kubecp/utils"
	"kubecp/utils/copyer"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MultiCopyQuery struct {
	Namespace string `json:"namespace" form:"namespace" binding:"required"`
	PodName   string `json:"pod_name" form:"pod_name" binding:"required"`
	DestPath  string `json:"dest_path" form:"dest_path" binding:"required"`
}

// @Summary MultiCopy2Container
// @description 上传到容器
// @Tags Kubernetes
// @Accept multipart/form-data
// @Param namespace query string true "namespace" default(default)
// @Param pod_name query string true "pod_name" default(nginx-test-76996486df-tdjdf)
// @Param dest_path query string false "dest_path" default(/root/)
// @Param files formData file true "files"
// @Success 200 {object} controller.JSONResult
// @Failure 500 {object} controller.JSONResult
// @Router /api/k8s/multi_upload [post]
func MultiCopy2Container(c *gin.Context) {
	render := controller.Gin{C: c}
	var query CopyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}
	// check namespace
	_, err := config.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), query.Namespace, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	podNameSlice := strings.Split(query.PodName, ",")
	for _, podName := range podNameSlice {
		// check pod
		pod, err := config.RestClient.CoreV1().Pods(query.Namespace).
			Get(context.TODO(), podName, metaV1.GetOptions{})
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
		if pod.Status.Phase == coreV1.PodSucceeded || pod.Status.Phase == coreV1.PodFailed {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase))
			return
		}
	}

	form, err := c.MultipartForm()
	if err != nil {
		render.SetError(utils.CODE_ERR_MSG, err)
		return
	}
	files := form.File["files"]
	var _tmpSaveDir = filepath.Join(config.TmpPath, strconv.FormatInt(time.Now().UnixNano(), 10))
	var fWg sync.WaitGroup
	var fErrCh = make(chan error, len(files))
	var fErr error
	go func() {
		for e := range fErrCh {
			fErr = multierror.Append(fErr, e)
		}
	}()

	// create _tmpSaveDir
	if !utils.FileOrPathExist(_tmpSaveDir) {
		err := os.MkdirAll(_tmpSaveDir, os.ModePerm)
		if err != nil {
			logs.Fatal("无法创建文件夹", err)
			render.SetError(utils.CODE_ERR_MSG, err)
			return
		}
	}
	defer os.RemoveAll(_tmpSaveDir)

	for _, file := range files {
		fWg.Add(1)
		go func(wg *sync.WaitGroup, fErrCh chan error, file *multipart.FileHeader) {
			defer wg.Done()
			// Default save path
			_tp := _tmpSaveDir
			uploadFilname := filepath.Base(file.Filename)
			uploadFPath := filepath.Dir(file.Filename)
			// Process folder upload
			if uploadFPath != "." {
				_tp = filepath.Join(_tmpSaveDir, uploadFPath)
				if !utils.FileOrPathExist(_tp) {
					os.MkdirAll(_tp, os.ModePerm)
				}
			}

			// save file to local in _tp
			err = c.SaveUploadedFile(file, filepath.Join(_tp, uploadFilname))
			if err != nil {
				fErrCh <- fmt.Errorf(file.Filename, err.Error())
			}
		}(&fWg, fErrCh, file)
	}
	fWg.Wait()

	time.Sleep(1 * time.Second)
	if fErr != nil {
		logs.Error(fErr.Error())
		render.SetError(utils.CODE_ERR_MSG, fErr)
		return
	}

	// strip trailing slash (if any)
	if query.DestPath != "/" && strings.HasSuffix(string(query.DestPath[len(query.DestPath)-1]), "/") {
		query.DestPath = query.DestPath[:len(query.DestPath)-1]
	}

	var res []string
	var cStopCh = make(chan struct{}, 1)
	var cErrCh = make(chan error, 1024)
	var copiedCh = make(chan string, 1024)
	var cErr []error
	go func() {
		for {
			select {
			case e := <-cErrCh:
				cErr = append(cErr, e)
			case copied := <-copiedCh:
				res = append(res, copied)
			case <-cStopCh:
				return
			}
		}
	}()
	var cWg sync.WaitGroup
	for _, podName := range podNameSlice {
		cWg.Add(1)
		var containerSlice []string
		res, err := config.RestClient.CoreV1().Pods(query.Namespace).
			Get(context.TODO(), podName, metaV1.GetOptions{})
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
		for _, container := range res.Spec.Containers {
			containerSlice = append(containerSlice, container.Name)
		}

		for _, container := range containerSlice {
			cWg.Add(1)
			go func(wg *sync.WaitGroup, cErrCh chan error, podName, container string) {
				defer wg.Done()
				reader, writer := io.Pipe()
				cp := copyer.NewCopyer(query.Namespace, podName, container, config.KuBeResConf, config.RestClient)
				cp.Stdin = reader

				go func() {
					defer writer.Close()
					err := utils.MakeTar(_tmpSaveDir, query.DestPath, writer)
					if err != nil {
						cErrCh <- fmt.Errorf("%s %v", container, err)
					}
				}()

				err := cp.CopyToPod(query.DestPath)
				if err != nil {
					cErrCh <- fmt.Errorf("pod: %s container: %s %v", podName, container, err)
					return
				}
				copiedCh <- fmt.Sprintf("pod: %s container: %s Copied", podName, container)
			}(&cWg, cErrCh, podName, container)
		}
		cWg.Done()
	}
	cWg.Wait()
	time.Sleep(1 * time.Second)

	if cErr != nil {
		var _se []string
		for _, e := range cErr {
			_se = append(_se, e.Error())
		}
		render.SetJson(map[string]interface{}{
			"success": fmt.Sprint(strings.Join(res, "<br>")),
			"failure": fmt.Sprint(strings.Join(_se, "<br>")),
		})
		return
	}
	render.SetJson(map[string]interface{}{
		"success": fmt.Sprint(strings.Join(res, "<br>")),
	})
}

type CopyQuery struct {
	Namespace     string `json:"namespace" form:"namespace" binding:"required"`
	PodName       string `json:"pod_name" form:"pod_name" binding:"required"`
	ContainerName string `json:"container_name" form:"container_name"`
	DestPath      string `json:"dest_path" form:"dest_path" binding:"required"`
}

// @Summary Copy2Container
// @description 上传到容器
// @Tags Kubernetes
// @Accept multipart/form-data
// @Param namespace query string true "namespace" default(default)
// @Param pod_name query string true "pod_name" default(nginx-test-76996486df-tdjdf)
// @Param container_name query string true "container_name" default(nginx-0)
// @Param dest_path query string false "dest_path" default(/root/)
// @Param files formData file true "files"
// @Success 200 {object} controller.JSONResult
// @Failure 500 {object} controller.JSONResult
// @Router /api/k8s/upload [post]
func Copy2Container(c *gin.Context) {
	render := controller.Gin{C: c}
	var query CopyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}
	// check namespace
	_, err := config.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), query.Namespace, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	// check pod
	pod, err := config.RestClient.CoreV1().Pods(query.Namespace).
		Get(context.TODO(), query.PodName, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	if pod.Status.Phase == coreV1.PodSucceeded || pod.Status.Phase == coreV1.PodFailed {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		render.SetError(utils.CODE_ERR_MSG, err)
		return
	}
	files := form.File["files"]
	var _tmpSaveDir = filepath.Join(config.TmpPath, strconv.FormatInt(time.Now().UnixNano(), 10))
	var fWg sync.WaitGroup
	var fErrCh = make(chan error, len(files))
	defer close(fErrCh)
	var fErr error
	go func() {
		for e := range fErrCh {
			fErr = multierror.Append(fErr, e)
		}
	}()

	// create _tmpSaveDir
	if !utils.FileOrPathExist(_tmpSaveDir) {
		err := os.MkdirAll(_tmpSaveDir, os.ModePerm)
		if err != nil {
			logs.Fatal("无法创建文件夹", err)
			render.SetError(utils.CODE_ERR_MSG, err)
			return
		}
	}
	defer os.RemoveAll(_tmpSaveDir)

	for _, file := range files {
		fWg.Add(1)
		go func(wg *sync.WaitGroup, fErrCh chan error, file *multipart.FileHeader) {
			defer wg.Done()
			// Default save path
			_tp := _tmpSaveDir
			uploadFilname := filepath.Base(file.Filename)
			uploadFPath := filepath.Dir(file.Filename)
			// Process folder upload
			if uploadFPath != "." {
				_tp = filepath.Join(_tmpSaveDir, uploadFPath)
				if !utils.FileOrPathExist(_tp) {
					os.MkdirAll(_tp, os.ModePerm)
				}
			}

			// save file to local in _tp
			err = c.SaveUploadedFile(file, filepath.Join(_tp, uploadFilname))
			if err != nil {
				fErrCh <- fmt.Errorf(file.Filename, err.Error())
			}
		}(&fWg, fErrCh, file)
	}

	fWg.Wait()
	time.Sleep(1 * time.Second)

	if fErr != nil {
		logs.Error(fErr.Error())
		render.SetError(utils.CODE_ERR_MSG, fErr)
		return
	}

	// strip trailing slash (if any)
	if query.DestPath != "/" && strings.HasSuffix(string(query.DestPath[len(query.DestPath)-1]), "/") {
		query.DestPath = query.DestPath[:len(query.DestPath)-1]
	}

	var containerSlice []string
	if len(query.ContainerName) != 0 {
		containerSlice = strings.Split(query.ContainerName, ",")
	} else {
		res, err := config.RestClient.CoreV1().Pods(query.Namespace).
			Get(context.TODO(), query.PodName, metaV1.GetOptions{})
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
		for _, container := range res.Spec.Containers {
			containerSlice = append(containerSlice, container.Name)
		}
	}

	var cWg sync.WaitGroup
	var cErrCh = make(chan error)
	var copiedCh = make(chan string)
	var cStopCh = make(chan struct{}, 1)
	var cErr []error
	var res []string

	go func() {
		for {
			select {
			case e := <-cErrCh:
				cErr = append(cErr, e)
			case copied := <-copiedCh:
				res = append(res, copied)
			case <-cStopCh:
				return
			}
		}
	}()

	for _, container := range containerSlice {
		cWg.Add(1)
		go func(wg *sync.WaitGroup, container string) {
			defer wg.Done()
			reader, writer := io.Pipe()
			cp := copyer.NewCopyer(query.Namespace, query.PodName, container, config.KuBeResConf, config.RestClient)
			cp.Stdin = reader

			go func() {
				defer writer.Close()
				err := utils.MakeTar(_tmpSaveDir, query.DestPath, writer)
				if err != nil {
					cErrCh <- fmt.Errorf("%s %v", container, err)
				}
			}()

			err := cp.CopyToPod(query.DestPath)
			if err != nil {
				cErrCh <- fmt.Errorf("container: %s %v", container, err)
				return
			}

			copiedCh <- fmt.Sprintf("container: %s Copied", container)
		}(&cWg, container)
	}
	cWg.Wait()
	time.Sleep(1 * time.Second)
	if cErr != nil {
		var _se []string
		for _, e := range cErr {
			_se = append(_se, e.Error())
		}
		render.SetJson(map[string]interface{}{
			"success": fmt.Sprint(strings.Join(res, "<br>")),
			"failure": fmt.Sprint(strings.Join(_se, "<br>")),
		})
		return
	}
	render.SetJson(map[string]interface{}{
		"success": fmt.Sprint(strings.Join(res, "<br>")),
	})
}

// @Summary Copy2Local
// @description 从容器下载到本地
// @Tags Kubernetes
// @Accept json
// @Param namespace query string true "namespace" default(default)
// @Param pod_name query string true "pod_name" default(nginx-test-76996486df-tdjdf)
// @Param container_name query string true "container_name" default(nginx-0)
// @Param dest_path query string false "dest_path" default(/root)
// @Success 200 {object} controller.JSONResult
// @Failure 500 {object} controller.JSONResult
// @Router /api/k8s/download [get]
func Copy2Local(c *gin.Context) {
	render := controller.Gin{C: c}
	var query CopyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	if len(query.ContainerName) == 0 {
		render.SetError(utils.CODE_ERR_APP, fmt.Errorf("ContainerName cannot be empty"))
		return
	}
	// check namespace
	_, err := config.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), query.Namespace, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	// check pod
	pod, err := config.RestClient.CoreV1().Pods(query.Namespace).
		Get(context.TODO(), query.PodName, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	if pod.Status.Phase == coreV1.PodSucceeded || pod.Status.Phase == coreV1.PodFailed {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase))
		return
	}

	reader, writer := io.Pipe()
	prefix := getPrefix(query.DestPath)
	prefix = path.Clean(prefix)
	// remove extraneous path shortcuts - these could occur if a path contained extra "../"
	// and attempted to navigate beyond "/" in a remote filesystem
	prefix = stripPathShortcuts(prefix)
	tarFileName := fmt.Sprintf("%s_%s.tar", strconv.FormatInt(time.Now().UnixNano(), 10), strings.ReplaceAll(prefix, "/", "_"))
	filePath := filepath.Join(config.TmpPath, tarFileName)

	//if err := utils.UnTarAll(r, _tmpSaveDir, prefix); err != nil {
	//	logs.Error(err)
	//	render.SetError(utils.CODE_ERR_APP, err)
	//	return
	//}

	fw, err := os.Create(filePath)
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	defer fw.Close()
	go func() {
		_, err = io.Copy(fw, reader)
		if err != nil {
			logs.Error(err)
		}
	}()

	cp := copyer.NewCopyer(query.Namespace, query.PodName, query.ContainerName, config.KuBeResConf, config.RestClient)
	cp.Stdout = writer

	err = cp.CopyFromPod(query.DestPath)
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	c.Header("X-Redirect", fmt.Sprintf("/%s/%s", config.Config.DownLoadTmp, tarFileName))
	render.SetJson(fmt.Sprintf("/%s/%s", config.Config.DownLoadTmp, tarFileName))
}

func getPrefix(file string) string {
	// tar strips the leading '/' if it's there, so we will too
	return strings.TrimLeft(file, "/")
}

// stripPathShortcuts removes any leading or trailing "../" from a given path
func stripPathShortcuts(p string) string {
	newPath := path.Clean(p)
	trimmed := strings.TrimPrefix(newPath, "../")

	for trimmed != newPath {
		newPath = trimmed
		trimmed = strings.TrimPrefix(newPath, "../")
	}

	// trim leftover {".", ".."}
	if newPath == "." || newPath == ".." {
		newPath = ""
	}

	if len(newPath) > 0 && string(newPath[0]) == "/" {
		return newPath[1:]
	}

	return newPath
}
