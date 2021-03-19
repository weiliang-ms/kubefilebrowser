package copyer

import (
	"fmt"
	"kubecp/logs"
	"path"
)

// 从io.Writer拷贝到pod
// tar file ---> io.Writer ---> kubernetes api ---> io.Reader ---> unTar file
func (c *copyer) CopyToPod(dest string) error {
	if c.NoPreserve {
		c.Command = []string{"tar", "--no-same-permissions", "--no-same-owner", "-xmf", "-"}
	} else {
		c.Command = []string{"tar", "-xmpf", "-"}
	}
	destDir := path.Dir(dest)
	if len(destDir) > 0 {
		c.Command = append(c.Command, "-C", destDir)
	}

	stderr, err := c.Exec()
	if err != nil {
		return fmt.Errorf(err.Error(), string(stderr))
	}
	logs.Warn(string(stderr))
	return nil
}
