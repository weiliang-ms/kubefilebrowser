package main

import (
	"fmt"
	"kubefilebrowser/cmd/kf_tools/command"
)

var (
	BuildAt string
	GitHash string
)

func main() {
	command.Execute(fmt.Sprintf("Hash: %s\nBuildDate: %s", GitHash, BuildAt))
}
