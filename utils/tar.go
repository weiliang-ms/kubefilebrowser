package utils

import (
	"archive/tar"
	"embed"
	_ "embed"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"kubefilebrowser/utils/logs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func MakeTar(srcPath, destPath string, writer io.Writer) error {
	// TODO: use compression here?
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	srcPath = path.Clean(srcPath)
	destPath = path.Clean(destPath)
	return recursiveTar(path.Dir(srcPath), path.Base(srcPath), path.Dir(destPath), path.Base(destPath), tarWriter)
}

func recursiveTar(srcBase, srcFile, destBase, destFile string, tw *tar.Writer) error {
	srcPath := path.Join(srcBase, srcFile)
	matchedPaths, err := filepath.Glob(srcPath)
	if err != nil {
		return err
	}
	for _, fPath := range matchedPaths {
		stat, err := os.Lstat(fPath)
		if err != nil {
			return err
		}
		if stat.IsDir() {
			files, err := ioutil.ReadDir(fPath)
			if err != nil {
				return err
			}
			if len(files) == 0 {
				//case empty directory
				hdr, _ := tar.FileInfoHeader(stat, fPath)
				hdr.Name = destFile
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}
			}
			for _, f := range files {
				if err := recursiveTar(srcBase, path.Join(srcFile, f.Name()), destBase, path.Join(destFile, f.Name()), tw); err != nil {
					return err
				}
			}
			return nil
		} else if stat.Mode()&os.ModeSymlink != 0 {
			//case soft link
			hdr, _ := tar.FileInfoHeader(stat, fPath)
			target, err := os.Readlink(fPath)
			if err != nil {
				return err
			}

			hdr.Linkname = target
			hdr.Name = destFile
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
		} else {
			//case regular file or other file type like pipe
			hdr, err := tar.FileInfoHeader(stat, fPath)
			if err != nil {
				return err
			}
			hdr.Name = destFile

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			f, err := os.Open(fPath)
			if err != nil {
				return err
			}

			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
			return f.Close()
		}
	}
	return nil
}

func UnTarAll(reader io.Reader, destDir, prefix string) error {
	tarReader := tar.NewReader(reader)
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		if !strings.HasPrefix(header.Name, prefix) {
			return fmt.Errorf("tar contents corrupted")
		}

		mode := header.FileInfo().Mode()
		destFileName := filepath.Join(destDir, header.Name[len(prefix):])

		baseName := filepath.Dir(destFileName)
		if err := os.MkdirAll(baseName, 0755); err != nil {
			return err
		}
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(destFileName, 0755); err != nil {
				return err
			}
			continue
		}

		evaledPath, err := filepath.EvalSymlinks(baseName)
		if err != nil {
			return err
		}

		if mode&os.ModeSymlink != 0 {
			linkname := header.Linkname

			if !filepath.IsAbs(linkname) {
				_ = filepath.Join(evaledPath, linkname)
			}

			if err := os.Symlink(linkname, destFileName); err != nil {
				return err
			}
		} else {
			outFile, err := os.Create(destFileName)
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			if err := outFile.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func isDestRelative(base, dest string) bool {
	relative, err := filepath.Rel(base, dest)
	if err != nil {
		return false
	}
	return relative == "." || relative == stripPathShortcuts(relative)
}

// stripPathShortcuts removes any leading or trailing "../" from a given path
func stripPathShortcuts(p string) string {
	newPath := path.Clean(p)
	trimmed := strings.TrimPrefix(newPath, "../")

	for trimmed != newPath {
		newPath = trimmed
		trimmed = strings.TrimPrefix(newPath, "../")
	}

	// trim leftover {".", ".."}
	if newPath == "." || newPath == ".." {
		newPath = ""
	}

	if len(newPath) > 0 && string(newPath[0]) == "/" {
		return newPath[1:]
	}

	return newPath
}

//go:embed ls_binary
var lsBinaryEmbededFiles embed.FS

func TarLs(lsPath string, writer io.Writer) error {
	fsys, err := fs.Sub(lsBinaryEmbededFiles, "ls_binary")
	if err != nil {
		logs.Error(err)
		return err
	}
	tw := tar.NewWriter(writer)
	// 如果关闭失败会造成tar包不完整
	defer tw.Close()
	f, err := fsys.Open(path.Base(lsPath))
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
	hdr, err := tar.FileInfoHeader(fi, lsPath)
	if err != nil {
		logs.Error(err)
		return err
	}
	hdr.Name = "ls"
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


//go:embed zip_binary
var zipBinaryEmbededFiles embed.FS

func TarZip(lsPath string, writer io.Writer) error {
	fsys, err := fs.Sub(zipBinaryEmbededFiles, "zip_binary")
	if err != nil {
		logs.Error(err)
		return err
	}
	tw := tar.NewWriter(writer)
	// 如果关闭失败会造成tar包不完整
	defer tw.Close()
	f, err := fsys.Open(path.Base(lsPath))
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
	hdr, err := tar.FileInfoHeader(fi, lsPath)
	if err != nil {
		logs.Error(err)
		return err
	}
	hdr.Name = "zip"
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