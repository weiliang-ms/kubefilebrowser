package utils

import (
	"archive/tar"
	"embed"
	_ "embed"
	"io"
	"io/fs"
	"kubefilebrowser/utils/logs"
	"path"
)

//go:embed kf_tools_binary
var kfBinaryEmbedFiles embed.FS

func TarKFTools(p string, writer io.Writer) error {
	fSys, err := fs.Sub(kfBinaryEmbedFiles, "kf_tools_binary")
	if err != nil {
		logs.Error(err)
		return err
	}
	tw := tar.NewWriter(writer)
	// 如果关闭失败会造成tar包不完整
	defer tw.Close()
	f, err := fSys.Open(path.Base(p))
	if err != nil {
		logs.Error(err)
		return err
	}
	f.Close()
	fi, err := f.Stat()
	if err != nil {
		logs.Error(err)
		return err
	}
	hdr, err := tar.FileInfoHeader(fi, p)
	if err != nil {
		logs.Error(err)
		return err
	}
	hdr.Name = "kf_tools"
	// 将tar的文件信息hdr写入到tw
	err = tw.WriteHeader(hdr)
	if err != nil {
		logs.Error(err)
		return err
	}
	// 将文件数据写入
	_, err = io.Copy(tw, f)

	return err
}
