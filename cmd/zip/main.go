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
    zw := zip.NewWriter(os.Stdout)
    //zw := zip.NewWriter(fw)
    defer func() {
        // 检测一下是否成功关闭
        if err := zw.Close(); err != nil {
            _, _ = fmt.Fprint(os.Stderr, err.Error())
        }
    }()
    var wg sync.WaitGroup
    for _, src := range paths {
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
        _, _ = fmt.Fprint(os.Stderr, err.Error())
    }
    
    for _, file := range files {
        if !file.IsDir() {
            dat, err := os.Open(inFilePath + "/" + file.Name())
            if err != nil {
                _, _ = fmt.Fprint(os.Stderr, err.Error())
            }
            // Add some files to the archive.
            f, err := zw.Create(inFilePath + "/" + file.Name())
            if err != nil {
                _, _ = fmt.Fprint(os.Stderr, err.Error())
            }
            _, err = io.Copy(f, dat)
            //_, err = f.Write(dat)
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
