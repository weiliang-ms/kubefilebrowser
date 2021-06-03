package terminalparser

import (
	"kubefilebrowser/utils/logs"
	"strings"
)

func ParseCmdOutput(p []byte) string {
	o := parse(p)
	return strings.Join(o, "\r\n")
}

func ParseCmdInput(p []byte) string {
	o := parse(p)
	return strings.Join(o, "")
}

func GetPs1(p []byte) string {
	lines := parse(p)
	if len(lines) == 0 {
		return ""
	}
	return lines[len(lines)-1]
}

func parse(p []byte) []string {
	defer func() {
		if r := recover(); r != nil {
			logs.Error(r)
		}
	}()
	s := Screen{
		Rows:   make([]*Row, 0, 1024),
		Cursor: &Cursor{},
	}
	return s.Parse(p)
}
