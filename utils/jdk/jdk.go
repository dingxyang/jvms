package jdk

import (
	"fmt"
	"github.com/tea4go/jvms/utils/file"
	"os"
)

// GetInstalled 获取已安装的 JDK 版本列表
// 参数:
//   root - JDK 安装的根目录路径
// 返回值:
//   []string - 已安装的 JDK 版本名称列表（按倒序排列）
func GetInstalled(root string) []string {
	list := make([]string, 0)
	files, _ := os.ReadDir(root)
	for i := len(files) - 1; i >= 0; i-- {
		if files[i].IsDir() {
			list = append(list, files[i].Name())
		}
	}
	return list
}

// IsVersionInstalled 检查指定版本的 JDK 是否已安装
// 参数:
//   root - JDK 安装的根目录路径
//   version - 要检查的 JDK 版本号
// 返回值:
//   bool - 已安装返回 true，否则返回 false
func IsVersionInstalled(root string, version string) bool {
	isInstalled := file.Exists(fmt.Sprintf("%s/%s/bin/javac.exe", root, version))
	return isInstalled
}
