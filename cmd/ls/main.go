package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type File struct {
	Name    string      `json:"Name"`
	Path    string      `json:"Path"`
	Size    int64       `json:"Size"`
	Mode    string      `json:"Mode"`
	ModTime time.Time   `json:"ModTime"`
	IsDir   bool        `json:"IsDir"`
	Sys     interface{} `json:"SysInfo"`
}

var path string

func init() {
	if len(os.Args) != 2 {
		path = "."
	} else {
		path = os.Args[1]
	}
}

var denyFileOrList = []string{"/ls", "/ls.exe", "/zip", "/zip.exe"}

func main() {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(254)
	}
	path = strings.Replace(path, "\\", "/", -1)
	path = strings.TrimRight(path, "/")
	var files []File
	for _, d := range dir {
		p := fmt.Sprintf("%s/%s", path, d.Name())
		if inSlice(p, denyFileOrList) {
			continue
		}
		isDir := d.IsDir()
		if d.Mode()&os.ModeSymlink != 0 {
			isDir = true
		}
		f := File{
			Name:    d.Name(),
			Path:    p,
			Size:    d.Size(),
			Mode:    d.Mode().String(),
			ModTime: d.ModTime(),
			IsDir:   isDir,
			Sys:     d.Sys(),
		}
		files = append(files, f)
	}
	s, err := json.Marshal(files)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(253)
	}
	fmt.Println(string(s))
}

func inSlice(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
