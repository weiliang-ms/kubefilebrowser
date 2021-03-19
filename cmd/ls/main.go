package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type File struct {
	Name    string      `json:"Name"`
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

func main() {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(254)
	}
	var files []File
	for _, d := range dir {
		if d.Name() == "ls" || d.Name() == "ls.exe" {
			continue
		}
		files = append(files, File{
			Name:    d.Name(),
			Size:    d.Size(),
			Mode:    d.Mode().String(),
			ModTime: d.ModTime(),
			IsDir:   d.IsDir(),
			Sys:     d.Sys(),
		})
	}
	s, err := json.Marshal(files)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(253)
	}
	fmt.Println(string(s))
}
