package copyer

import (
	"fmt"
	"path"
	"time"
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
	// 重试三次
	attempts := 3
	attempt := 0
	for attempt < attempts {
		attempt++
		stderr, err := c.Exec()
		if len(stderr) != 0 {
			if attempt == attempts {
				return fmt.Errorf("STDERR: " + string(stderr))
			}
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		if err != nil {
			if attempt == attempts {
				return err
			}
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		return nil
	}
	return nil
}
