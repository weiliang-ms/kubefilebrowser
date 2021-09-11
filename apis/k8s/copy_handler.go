package k8s

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubefilebrowser/apis"
	"kubefilebrowser/configs"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/copyer"
	"kubefilebrowser/utils/execer"
	"kubefilebrowser/utils/logs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MultiCopyQuery struct {
	Namespace string   `json:"namespace" form:"namespace" binding:"required"`
	PodName   []string `json:"pod_name" form:"pod_name" binding:"required"`
	DestPath  string   `json:"dest_path" form:"dest_path" binding:"required"`
}

// MultiCopy2Container
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
	check := Check{
		namespace: query.Namespace,
	}
	// check namespace
	if _, err := check.Namespace(); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	for _, podName := range query.PodName {
		// check pod
		check.pod = podName
		_, err := check.Pod()
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
	}
	var _tmpSaveDir = filepath.Join(configs.TmpPath, strconv.FormatInt(time.Now().UnixNano(), 10))
	defer os.RemoveAll(_tmpSaveDir)
	if err := writeFiles(c, _tmpSaveDir); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_MSG, err)
		return
	}

	// strip trailing slash (if any)
	if query.DestPath != "/" && strings.HasSuffix(string(query.DestPath[len(query.DestPath)-1]), "/") {
		query.DestPath = query.DestPath[:len(query.DestPath)-1]
	}

	var res []string
	var mu sync.Mutex
	var wg sync.WaitGroup
	var err error
	var cErr []error
	for _, podName := range query.PodName {
		var containerSlice []string
		var pod *coreV1.Pod
		pod, err = configs.RestClient.CoreV1().Pods(query.Namespace).
			Get(context.TODO(), podName, metaV1.GetOptions{})
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
		for _, container := range pod.Spec.Containers {
			containerSlice = append(containerSlice, container.Name)
		}

		for _, container := range containerSlice {
			wg.Add(1)
			go func(wg *sync.WaitGroup, podName, container string) {
				defer wg.Done()
				err = copyToPod(query.Namespace, podName, container, _tmpSaveDir, query.DestPath)
				if err != nil {
					mu.Lock()
					cErr = append(cErr, err)
					mu.Unlock()
					return
				}
				mu.Lock()
				res = append(res, fmt.Sprintf("pod %s container %s copied", query.PodName, container))
				mu.Unlock()
			}(&wg, podName, container)
		}
	}
	wg.Wait()

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
	Namespace     string   `json:"namespace" form:"namespace" binding:"required"`
	PodName       string   `json:"pod_name" form:"pod_name" binding:"required"`
	ContainerName []string `json:"container_name" form:"container_name"`
	DestPath      string   `json:"dest_path" form:"dest_path" binding:"required"`
}

// Copy2Container
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
	check := Check{
		namespace: query.Namespace,
		pod:       query.PodName,
	}
	// check namespace
	if _, err := check.Namespace(); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	// check pod
	if _, err := check.Pod(); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	var _tmpSaveDir = filepath.Join(configs.TmpPath, strconv.FormatInt(time.Now().UnixNano(), 10))
	defer os.RemoveAll(_tmpSaveDir)
	if err := writeFiles(c, _tmpSaveDir); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_MSG, err)
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


	var cErr []error
	var res []string
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, container := range containerSlice {
		wg.Add(1)
		go func(wg *sync.WaitGroup, container string) {
			defer wg.Done()
			err := copyToPod(query.Namespace, query.PodName, container, _tmpSaveDir, query.DestPath)
			if err != nil {
				mu.Lock()
				cErr = append(cErr, err)
				mu.Unlock()
				return
			}
			mu.Lock()
			res = append(res, fmt.Sprintf("pod %s container %s copied", query.PodName, container))
			mu.Unlock()
		}(&wg, container)
	}
	wg.Wait()
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

func copyToPod(namespace, pod, container, _tmpSaveDir, destPath string) error {
	reader, writer := io.Pipe()
	cp := copyer.NewCopyer(namespace, pod, container, configs.KuBeResConf, configs.RestClient)
	cp.Stdin = reader

	go func() {
		defer writer.Close()
		err := utils.MakeTar(_tmpSaveDir, destPath, writer)
		if err != nil {
			logs.Error(err)
		}
	}()

	return cp.CopyToPod(destPath)
}

type CopyFromPodQuery struct {
	Namespace     string   `json:"namespace" form:"namespace" binding:"required"`
	PodName       string   `json:"pod_name" form:"pod_name" binding:"required"`
	ContainerName string   `json:"container_name" form:"container_name"`
	DestPath      []string `json:"dest_path" form:"dest_path" binding:"required"`
	Style         string   `json:"style" form:"style"`
}

// Copy2Local
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

	check := Check{
		namespace: query.Namespace,
		pod:       query.PodName,
	}
	// check namespace
	if _, err := check.Namespace(); err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	// check pod
	pod, err := check.Pod()
	if err != nil {
		logs.Error(err)
		render.SetError(utils.CODE_ERR_APP, err)
		return
	}

	//var isUnix = true
	//for _, value := range pod.Spec.NodeSelector {
	//	if strings.Contains(value, "windows") {
	//		isUnix = false
	//		break
	//	}
	//}
	var osType = "linux"
	var arch = "amd64"

	// get pod system arch and type
	node, err := configs.RestClient.CoreV1().Nodes().
		Get(context.TODO(), pod.Spec.NodeName, metaV1.GetOptions{})
	if err == nil {
		if node.Labels["beta.kubernetes.io/os"] != "" {
			osType = node.Labels["beta.kubernetes.io/os"]
		} else if node.Labels["kubernetes.io/os"] != "" {
			osType = node.Labels["kubernetes.io/os"]
		}
		if node.Labels["beta.kubernetes.io/arch"] != "" {
			arch = node.Labels["beta.kubernetes.io/arch"]
		} else if node.Labels["kubernetes.io/arch"] != "" {
			arch = node.Labels["kubernetes.io/arch"]
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
		zipPath := fmt.Sprintf("/kf_tools_%s_%s", osType, arch)
		if osType == "windows" {
			zipPath = fmt.Sprintf("/kf_tools_%s_%s.exe", osType, arch)
		}
		err = query.copyZipTar(zipPath)
		if err != nil {
			logs.Error(err)
			render.SetError(utils.CODE_ERR_APP, err)
			return
		}
		if osType != "windows" {
			_cmd := []string{"chmod", "+x", "/kf_tools"}
			_, err = query.exec(_cmd)
			if err != nil {
				logs.Error(err)
				render.SetError(utils.CODE_ERR_APP, err)
				return
			}
		}
		reader, writer := io.Pipe()
		cp := copyer.NewCopyer(query.Namespace, query.PodName, query.ContainerName, configs.KuBeResConf, configs.RestClient)
		cp.Stdin = reader
		cp.Stdout = render.C.Writer
		go func() {
			<-c.Request.Context().Done()
			_, _ = writer.Write([]byte("close\n"))
		}()
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
		err := utils.TarKFTools(zipPath, writer)
		if err != nil {
			logs.Error(err)
		}
	}()
	return cp.CopyToPod(zipPath)
}
