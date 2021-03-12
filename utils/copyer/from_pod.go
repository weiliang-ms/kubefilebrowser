package copyer

import (
	"fmt"
	"time"
)

// 从pod内拷贝到io.Writer
func (c *copyer) CopyFromPod(dest []string) error {
	c.Command = append([]string{"tar", "cf", "-"}, dest...)
	attempts := 3
	attempt := 0
	for attempt < attempts {
		attempt++

		stderr, err := c.Exec()
		if attempt == attempts {
			if len(stderr) != 0 {
				return fmt.Errorf("STDERR: " + (string)(stderr))
			}
			if err != nil {
				return err
			}
		}
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	return nil
}
