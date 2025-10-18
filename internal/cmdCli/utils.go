package cmdCli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tea4go/jvms/internal/entity"
	"github.com/tea4go/jvms/utils/jdk"
)

// getSimilarAvailableVersions 根据版本获取相似的可用版本，用于支持版本未找到错误提示
// func getSimilarAvailableVersions(version string) {

// }

// getJavaHome 从临时JDK文件目录中获取JAVA_HOME路径
// 参数:
//   jdkTempFile - JDK临时解压目录路径
// 返回值:
//   string - JAVA_HOME路径(包含javac.exe的父目录)
func getJavaHome(jdkTempFile string) string {
	var javaHome string
	fs.WalkDir(os.DirFS(jdkTempFile), ".", func(path string, d fs.DirEntry, err error) error {
		if filepath.Base(path) == "javac.exe" {
			temPath := strings.Replace(path, "bin/javac.exe", "", -1)
			javaHome = filepath.Join(jdkTempFile, temPath)
			return fs.SkipDir
		}
		return nil
	})
	return javaHome
}

// getJdkVersions 获取可供下载的JDK版本列表
// 参数:
//   cfx - 配置对象指针
// 返回值:
//   []entity.JdkVersion - JDK版本列表
//   error - 错误信息
func getJdkVersions(cfx *entity.Config) ([]entity.JdkVersion, error) {
	var versions []entity.JdkVersion

	fmt.Println("")
	fmt.Println("-= Huawei OpenJDK Mirror =-")
	// 华为镜像 JDKs
	huaweiJdks := jdk.HuaweiJDKs()
	for _, huaweiJdk := range huaweiJdks {
		versionName := fmt.Sprintf("openjdk-%s", huaweiJdk.Version)
		fmt.Printf("%s [%s/%s] %s\n", versionName, huaweiJdk.GOOS, huaweiJdk.GOARCH, huaweiJdk.URL)
		versions = append(versions, entity.JdkVersion{Version: versionName, Url: huaweiJdk.URL})
	}

	return versions, nil
}
