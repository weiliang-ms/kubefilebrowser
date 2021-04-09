package copyer

import (
	"fmt"
	"kubefilebrowser/utils/logs"
	"path"
	"strings"
)

// 从io.Writer拷贝到pod
// tar file ---> io.Writer ---> kubernetes api ---> io.Reader ---> unTar file
func (c *copyer) CopyToPod(dest string) error {
	if c.NoPreserve {
		c.Command = []string{"tar", "--no-same-permissions", "--no-same-owner", "-xmf", "-"}
	} else {
		c.Command = []string{"tar", "-xmf", "-"}
	}
	destDir := path.Dir(dest)
	if len(destDir) > 0 {
		c.Command = append(c.Command, "-C", destDir)
	}

	stderr, err := c.Exec()
	if err != nil {
		return fmt.Errorf(err.Error(), string(stderr))
	}
	if len(stderr) != 0 {
		logs.Warn(string(stderr))
		for _, line := range strings.Split(string(stderr), "\n") {
			if len(strings.TrimSpace(line)) == 0 {
				continue
			}
			if !strings.Contains(strings.ToLower(line), "removing") {
				return fmt.Errorf(line)
			}
		}
	}
	return nil
}
