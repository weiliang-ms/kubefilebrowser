package copyer

import "fmt"

// 从pod内拷贝到io.Writer
func (c *copyer) CopyFromPod(dest []string, style string) error {
	switch style {
	case "tar":
		c.Command = append([]string{"tar", "cf", "-"}, dest...)
	case "zip":
		c.Command = append([]string{"/kf_tools", "zip"}, dest...)
	default:
		c.Command = append([]string{"tar", "cf", "-"}, dest...)
	}
    stderr, err := c.Exec()
    if err != nil {
        if len(stderr) != 0 {
            return fmt.Errorf("STDERR: " + string(stderr))
        }
        return err
    }
    // 三次重试
	//attempts := 3
	//attempt := 0
	//for attempt < attempts {
	//	attempt++
    //
	//	stderr, err := c.Exec()
	//	logs.Error(err, string(stderr))
	//	if attempt == attempts {
	//		if err != nil {
	//			return err
	//		}
	//		if len(stderr) != 0 {
	//			return fmt.Errorf("STDERR: " + string(stderr))
	//		}
	//	}
	//	if err == nil {
	//		return nil
	//	}
	//	time.Sleep(time.Duration(attempt) * time.Second)
	//}
	return nil
}
