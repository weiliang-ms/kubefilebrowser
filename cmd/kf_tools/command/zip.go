package command

import (
	"archive/zip"
	"github.com/spf13/cobra"
	"io"
	"kubefilebrowser/utils"
	"kubefilebrowser/utils/ratelimit"
	"kubefilebrowser/utils/symwalk"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(zipCmd)
}

// zipCmd represents the zip command
var zipCmd = &cobra.Command{
	Use:   "zip",
	Short: "The default action is to add or replace zipfile entries from list",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		zw := zip.NewWriter(ratelimit.Writer(os.Stdout, ratelimit.New(2000*1024))) // 限制输出 2000KB/s
		defer zw.Close()

		for _, p := range args {
			if !utils.FileOrPathExist(p) {
				continue
			}
			makeZip(p, zw)
		}
	},
}

func makeZip(inFilepath string, zw *zip.Writer) error {
	return symwalk.Walk(inFilepath, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		// 目录拉平
		//relPath := strings.TrimPrefix(filePath, filepath.Dir(inFilepath))
		var zwPath = utils.ToLinuxPath(filePath)
		zipFile, err := zw.Create(strings.TrimPrefix(zwPath, "/"))
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fsFile.Close()
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
}
