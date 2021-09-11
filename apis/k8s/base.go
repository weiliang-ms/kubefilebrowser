package k8s

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-multierror"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubefilebrowser/configs"
	"kubefilebrowser/utils"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"
)

type Check struct {
	namespace string
	pod       string
}

func (c *Check) Namespace() (*coreV1.Namespace, error) {
	ns, err := configs.RestClient.CoreV1().Namespaces().
		Get(context.TODO(), c.namespace, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func (c *Check) Pod() (*coreV1.Pod, error) {
	pod, err := configs.RestClient.CoreV1().Pods(c.namespace).
		Get(context.TODO(), c.pod, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if pod.Status.Phase == coreV1.PodSucceeded || pod.Status.Phase == coreV1.PodFailed {
		return nil, fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase)
	}
	return pod, nil
}

func writeFiles(c *gin.Context, _tmpSaveDir string) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]
	if len(files) == 0 {
		return fmt.Errorf("files is null")
	}
	var wg sync.WaitGroup
	var fErr error

	// create _tmpSaveDir
	if !utils.FileOrPathExist(_tmpSaveDir) {
		err = os.MkdirAll(_tmpSaveDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	for _, f := range files {
		wg.Add(1)
		go func(wg *sync.WaitGroup, file *multipart.FileHeader) {
			defer wg.Done()
			// Default save path
			_tp := _tmpSaveDir
			uploadFileName := filepath.Base(file.Filename)
			uploadFPath := filepath.Dir(file.Filename)
			// Process folder upload
			if uploadFPath != "." {
				_tp = filepath.Join(_tmpSaveDir, uploadFPath)
				if !utils.FileOrPathExist(_tp) {
					_ = os.MkdirAll(_tp, os.ModePerm)
				}
			}

			// save file to local in _tp
			err = c.SaveUploadedFile(file, filepath.Join(_tp, uploadFileName))
			if err != nil {
				fErr = multierror.Append(fErr, fmt.Errorf(file.Filename, err.Error()))
			}
		}(&wg, f)
	}
	wg.Wait()
	if fErr != nil {
		return fErr
	}
	return nil
}
