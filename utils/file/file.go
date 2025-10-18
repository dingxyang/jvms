package file

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Unzip 解压缩 zip 文件到指定目录
// 参数:
//   src - zip 文件的源路径
//   dest - 解压缩的目标目录路径
// 返回值:
//   error - 解压过程中的错误，成功则返回 nil
// 函数来源: http://stackoverflow.com/users/1129149/swtdrgn
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Exists 检查文件或目录是否存在
// 参数:
//   filename - 要检查的文件或目录路径
// 返回值:
//   bool - 存在返回 true，否则返回 false
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// GetCurrentPath 获取当前可执行文件所在的目录路径
// 返回值:
//   string - 当前可执行文件的目录路径，获取失败则返回空字符串
func GetCurrentPath() string {
	currentDir, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(currentDir)
}
