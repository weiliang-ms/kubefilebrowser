package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var paths []string

func init() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "path is empty")
		os.Exit(250)
	}
	paths = os.Args[1:]
}

func main() {
	// debug code
	//fw, err := os.Create("xxxx.zip")
	//if err != nil {
	//	_, _ = fmt.Fprint(os.Stderr, err.Error())
	//	return
	//}
	//zw := zip.NewWriter(fw)

	zw := zip.NewWriter(os.Stdout)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			_, _ = fmt.Fprint(os.Stderr, err.Error())
		}
	}()
	var wg sync.WaitGroup
	for _, src := range paths {
		if !fileOrPathExist(src) {
			_, _ = fmt.Fprintf(os.Stdout, "%s does not exist", src)
			_, _ = fmt.Fprintf(os.Stderr, "%s does not exist", src)
			continue
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, src string, zw *zip.Writer) {
			defer wg.Done()
			makeZip(src, zw)

		}(&wg, src, zw)
	}
	wg.Wait()
}

// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
func makeZip(inFilePath string, zw *zip.Writer) {
	inFilePath = strings.Replace(inFilePath, "\\", "/", -1)
	_inFilePath := strings.Split(inFilePath, ":")
	if len(_inFilePath) >= 2 {
		inFilePath = strings.Join(_inFilePath[1:], ":")
	}
	files, err := ioutil.ReadDir(inFilePath)
	if err != nil {
		file, err := os.Stat(inFilePath)
		if err != nil {
			_, _ = fmt.Fprint(os.Stderr, err.Error())
		}
		files = append(files, file)
	}

	for _, file := range files {
		if !file.IsDir() {
			var dat *os.File
			var fw io.Writer
			var err error
			_p := strings.Split(inFilePath, "/")
			if _p[len(_p)-1] == file.Name() {
				dat, err = os.Open(inFilePath)
			} else {
				dat, err = os.Open(inFilePath + "/" + file.Name())
			}
			if err != nil {
				_, _ = fmt.Fprint(os.Stderr, err.Error())
				continue
			}
			// Add some files to the archive.
			if _p[len(_p)-1] == file.Name() {
				fw, err = zw.Create(strings.Replace(inFilePath, "/", "", 1))
			} else {
				fw, err = zw.Create(strings.Replace(inFilePath, "/", "", 1) + "/" + file.Name())
			}
			if err != nil {
				_, _ = fmt.Fprint(os.Stderr, err.Error())
				continue
			}
			_, err = io.Copy(fw, dat)
			if err != nil {
				_, _ = fmt.Fprint(os.Stderr, err.Error())
			}
		} else {
			// Recurse
			newBase := inFilePath + "/" + file.Name()
			makeZip(newBase, zw)
		}
	}
}

func fileOrPathExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
