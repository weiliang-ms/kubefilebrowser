package k8s

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	"io"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubefilebrowser/apis"
	"kubefilebrowser/configs"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/copyer"
	"kubefilebrowser/utils/execer"
	"kubefilebrowser/utils/logs"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MultiCopyQuery struct {
	Namespace string `json:"namespace" form:"namespace" binding:"required"`
	PodName   []string `json:"pod_name" form:"pod_name" binding:"required"`
	DestPath  string `json:"dest_path" form:"dest_path" binding:"required"`
}

// @Summary MultiCopy2Container
// @description 上传到容器
// @Tags Kubernetes
// @Accept multipart/form-data
// @Param namespace query MultiCopyQuery true "namespace" default(default)
// @Param pod_name query MultiCopyQuery true "pod_name" default(nginx-test-76996486df-tdjdf)
// @Param dest_path query MultiCopyQuery false "dest_path" default(/root/)
// @Param files formData file true "files"
// @Success 200 {object} apis.JSONResult
// @Failure 500 {object} apis.JSONResult
// @Router /api/k8s/multi_upload [post]
func MultiCopy2Container(c *gin.Context) {
	render := apis.Gin{C: c}
	var query MultiCopyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}
	// check namespace
	_, err := configs.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), query.Namespace, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}
	for _, podName := range query.PodName {
		// check pod
		pod, err := configs.RestClient.CoreV1().Pods(query.Namespace).
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
	if len(files) == 0 {
		logs.Error("files is null")
		render.SetError(utils.CODE_ERR_MSG, fmt.Errorf("files is null"))
		return
	}

	var _tmpSaveDir = filepath.Join(configs.TmpPath, strconv.FormatInt(time.Now().UnixNano(), 10))
	var wg sync.WaitGroup
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
		wg.Add(1)
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
		}(&wg, fErrCh, file)
	}
	wg.Wait()

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

	for _, podName := range query.PodName {
		var containerSlice []string
		res, err := configs.RestClient.CoreV1().Pods(query.Namespace).
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
			wg.Add(1)
			go func(wg *sync.WaitGroup, cErrCh chan error, podName, container string) {
				defer wg.Done()
				reader, writer := io.Pipe()
				cp := copyer.NewCopyer(query.Namespace, podName, container, configs.KuBeResConf, configs.RestClient)
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
					logs.Error(fmt.Sprintf("pod: %s container: %s %v", podName, container, err))
					cErrCh <- fmt.Errorf("pod: %s container: %s %v", podName, container, err)
					return
				}
				logs.Info(fmt.Sprintf("pod: %s container: %s Copied", podName, container))
				copiedCh <- fmt.Sprintf("pod: %s container: %s Copied", podName, container)
			}(&wg, cErrCh, podName, container)
		}
	}
	wg.Wait()
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
	ContainerName []string `json:"container_name" form:"container_name"`
	DestPath      string `json:"dest_path" form:"dest_path" binding:"required"`
}

// @Summary Copy2Container
// @description 上传到容器
// @Tags Kubernetes
// @Accept multipart/form-data
// @Param namespace query CopyQuery true "namespace" default(default)
// @Param pod_name query CopyQuery true "pod_name" default(nginx-test-76996486df-tdjdf)
// @Param container_name query CopyQuery true "container_name" default(nginx-0)
// @Param dest_path query CopyQuery false "dest_path" default(/root/)
// @Param files formData file true "files"
// @Success 200 {object} apis.JSONResult
// @Failure 500 {object} apis.JSONResult
// @Router /api/k8s/upload [post]
func Copy2Container(c *gin.Context) {
	render := apis.Gin{C: c}
	var query CopyQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_PARAM, err)
		return
	}
	// check namespace
	_, err := configs.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), query.Namespace, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	// check pod
	pod, err := configs.RestClient.CoreV1().Pods(query.Namespace).
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
		logs.Error(err)
		render.SetError(utils.CODE_ERR_MSG, err)
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		logs.Error("files is null")
		render.SetError(utils.CODE_ERR_MSG, fmt.Errorf("files is null"))
		return
	}

	var _tmpSaveDir = filepath.Join(configs.TmpPath, strconv.FormatInt(time.Now().UnixNano(), 10))
	var wg sync.WaitGroup
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
	//defer os.RemoveAll(_tmpSaveDir)

	for _, file := range files {
		wg.Add(1)
		go func(wg *sync.WaitGroup, fErrCh chan error, file *multipart.FileHeader) {
			defer wg.Done()
			// Default save path
			logs.Debug(file.Filename)
			_tp := _tmpSaveDir
			uploadFilname := filepath.Base(file.Filename)
			uploadFPath := filepath.Dir(file.Filename)
			// Process folder upload
			if uploadFPath != "." {
				_tp = filepath.Join(_tmpSaveDir, uploadFPath)
				if !utils.FileOrPathExist(_tp) {
					err = os.MkdirAll(_tp, os.ModePerm)
					if err != nil {
						logs.Error(err)
					}
				}
			}

			// save file to local in _tp
			err = c.SaveUploadedFile(file, filepath.Join(_tp, uploadFilname))
			if err != nil {
				fErrCh <- fmt.Errorf(file.Filename, err.Error())
			}
		}(&wg, fErrCh, file)
	}

	wg.Wait()
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
		containerSlice = query.ContainerName
	} else {
		res, err := configs.RestClient.CoreV1().Pods(query.Namespace).
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
		wg.Add(1)
		go func(wg *sync.WaitGroup, container string) {
			defer wg.Done()
			reader, writer := io.Pipe()
			cp := copyer.NewCopyer(query.Namespace, query.PodName, container, configs.KuBeResConf, configs.RestClient)
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
				logs.Error("container: ", container, err)
				cErrCh <- fmt.Errorf("container: %s %v", container, err)
				return
			}
			logs.Info("container: ", container, " Copied")
			copiedCh <- fmt.Sprintf("container: %s Copied", container)
		}(&wg, container)
	}
	wg.Wait()
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

type CopyFromPodQuery struct {
	Namespace     string   `json:"namespace" form:"namespace" binding:"required"`
	PodName       string   `json:"pod_name" form:"pod_name" binding:"required"`
	ContainerName string   `json:"container_name" form:"container_name"`
	DestPath      []string `json:"dest_path" form:"dest_path" binding:"required"`
	Style         string   `json:"style" form:"style"`
}

// @Summary Copy2Local
// @description 从容器下载到本地
// @Tags Kubernetes
// @Accept json
// @Param namespace query CopyFromPodQuery true "namespace" default(default)
// @Param pod_name query CopyFromPodQuery true "pod_name" default(nginx-test-76996486df-tdjdf)
// @Param container_name query CopyFromPodQuery true "container_name" default(nginx-0)
// @Param dest_path query CopyFromPodQuery true "dest_path" default(/root)
// @Param style query CopyFromPodQuery true "style" default(rar)
// @Success 200 {object} apis.JSONResult
// @Failure 500 {object} apis.JSONResult
// @Router /api/k8s/download [get]
func Copy2Local(c *gin.Context) {
	render := apis.Gin{C: c}
	var query = CopyFromPodQuery{
		Style: "rar",
	}
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
	_, err := configs.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), query.Namespace, metaV1.GetOptions{})
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	// check pod
	pod, err := configs.RestClient.CoreV1().Pods(query.Namespace).
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
	
	var isUnix = true
	for _, value := range pod.Spec.NodeSelector {
		if strings.Contains(value, "windows") {
			isUnix = false
			break
		}
	}
	
	fileName := fmt.Sprintf("%s.%s", strconv.FormatInt(time.Now().UnixNano(), 10), query.Style)
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("X-File-Name", fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	switch query.Style {
	case "tar":
		cp := copyer.NewCopyer(query.Namespace, query.PodName, query.ContainerName, configs.KuBeResConf, configs.RestClient)
		cp.Stdout = render.C.Writer
		err = cp.CopyFromPod(query.DestPath, query.Style)
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
	case "zip":
		zipPath := "/zip_linux_amd64"
		if !isUnix {
			zipPath = "/zip_windows_amd64.exe"
		}
		err = query.copyZipTar(zipPath)
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
		if isUnix {
			_cmd := []string{"chmod", "+x", "/zip"}
			_, err = query.exec(_cmd)
			if err != nil {
				logs.Error(err)
				render.SetError(utils.CODE_ERR_APP, err)
				return
			}
		}
		
		cp := copyer.NewCopyer(query.Namespace, query.PodName, query.ContainerName, configs.KuBeResConf, configs.RestClient)
		cp.Stdout = render.C.Writer
		err = cp.CopyFromPod(query.DestPath, query.Style)
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
	default:
		render.SetError(utils.CODE_ERR_MSG, fmt.Errorf("no matching compression type found"))
		return
	}
}

func (query *CopyFromPodQuery) exec(command []string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	exec := execer.NewExec(query.Namespace, query.PodName, query.ContainerName, configs.KuBeResConf, configs.RestClient)
	exec.Command = command
	exec.Tty = false
	exec.Stdout = &stdout
	exec.Stderr = &stderr
	err := exec.Exec()
	if err != nil {
		logs.Error(stderr.String())
		return nil, err
	}
	if len(stderr.Bytes()) != 0 {
		return nil, fmt.Errorf(stderr.String())
	}
	
	return stdout.Bytes(), nil
}

func (query *CopyFromPodQuery) copyZipTar(zipPath string) error {
	reader, writer := io.Pipe()
	cp := copyer.NewCopyer(query.Namespace, query.PodName, query.ContainerName, configs.KuBeResConf, configs.RestClient)
	cp.Stdin = reader
	
	go func() {
		defer writer.Close()
		err := utils.TarZip(zipPath, writer)
		if err != nil {
			logs.Error(err)
		}
	}()
	return cp.CopyToPod(zipPath)
}