package utils

import "time"

type File struct {
    Name    string      `json:"Name"`
    Size    int64       `json:"Size"`
    Mode    string      `json:"Mode"`
    ModTime time.Time   `json:"ModTime"`
    IsDir   bool        `json:"IsDir"`
    Sys     interface{} `json:"SysInfo"`
}
