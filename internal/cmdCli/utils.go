package cmdCli

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/tea4go/jvms/internal/entity"
	"github.com/tea4go/jvms/utils/jdk"
)

// getSimilarAvailableVersions 根据版本获取相似的可用版本，用于支持版本未找到错误提示
// func getSimilarAvailableVersions(version string) {

// }

// getJavaHome 从临时JDK文件目录中获取JAVA_HOME路径
// 参数:
//
//	jdkTempFile - JDK临时解压目录路径
//
// 返回值:
//
//	string - JAVA_HOME路径(包含javac.exe的父目录)
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

// parseVersion 解析版本号字符串为数字切片
// 参数:
//
//	version - 版本号字符串，例如 "openjdk-12.0.1" 或 "12.0.1"
//
// 返回值:
//
//	[]int - 版本号数字切片，例如 [12, 0, 1]
func parseVersion(version string) []int {
	// 移除 "openjdk-" 前缀
	version = strings.TrimPrefix(version, "openjdk-")

	// 按点分割版本号
	parts := strings.Split(version, ".")
	numbers := make([]int, 0, len(parts))

	for _, part := range parts {
		// 转换为整数
		if num, err := strconv.Atoi(part); err == nil {
			numbers = append(numbers, num)
		} else {
			// 如果转换失败，视为0
			numbers = append(numbers, 0)
		}
	}

	return numbers
}

// compareVersions 比较两个版本号
// 参数:
//
//	v1 - 版本1
//	v2 - 版本2
//
// 返回值:
//
//	int - 如果 v1 < v2 返回 -1，v1 == v2 返回 0，v1 > v2 返回 1
func compareVersions(v1, v2 string) int {
	nums1 := parseVersion(v1)
	nums2 := parseVersion(v2)

	// 比较每个数字段
	maxLen := len(nums1)
	if len(nums2) > maxLen {
		maxLen = len(nums2)
	}

	for i := 0; i < maxLen; i++ {
		n1 := 0
		n2 := 0

		if i < len(nums1) {
			n1 = nums1[i]
		}
		if i < len(nums2) {
			n2 = nums2[i]
		}

		if n1 < n2 {
			return -1
		} else if n1 > n2 {
			return 1
		}
	}

	return 0
}

// extractMajorVersion 提取主版本号
// 参数:
//
//	version - 完整版本号字符串，例如 "openjdk-12.0.1"
//
// 返回值:
//
//	string - 主版本号，例如 "12"
func extractMajorVersion(version string) string {
	// 移除 "openjdk-" 前缀
	version = strings.TrimPrefix(version, "openjdk-")

	// 按点分割，取第一段
	parts := strings.Split(version, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return version
}

// getJdkVersions 获取可供下载的JDK版本列表
// 参数:
//
//	cfx - 配置对象指针
//
// 返回值:
//
//	[]entity.TJDKVersion - JDK版本列表（按版本号从大到小排序，每个主版本号只显示最新的完整版本）
//	error - 错误信息
func getJdkVersions(cfx *entity.TConfig) ([]entity.TJDKVersion, error) {
	var versions []entity.TJDKVersion

	//fmt.Println("")
	//fmt.Println("-= Huawei OpenJDK Mirror =-")
	// 华为镜像 JDKs
	huaweiJdks := jdk.HuaweiJDKs()

	// 使用 map 来去重，只保留每个主版本号的最新完整版本
	majorVersionMap := make(map[string]entity.TJDKVersion)

	for _, huaweiJdk := range huaweiJdks {
		versionName := fmt.Sprintf("openjdk-%s", huaweiJdk.Version)
		majorVersion := extractMajorVersion(versionName)

		// 如果该主版本号还没有记录，或者当前版本更新，则保存完整版本号
		if existing, exists := majorVersionMap[majorVersion]; !exists {
			majorVersionMap[majorVersion] = entity.TJDKVersion{
				Version: versionName, // 保留完整版本号，例如 "openjdk-12.0.2"
				Url:     huaweiJdk.URL,
			}
		} else {
			// 比较版本，保留更新的版本（完整版本号）
			if compareVersions(versionName, existing.Version) > 0 {
				majorVersionMap[majorVersion] = entity.TJDKVersion{
					Version: versionName, // 保留完整版本号
					Url:     huaweiJdk.URL,
				}
			}
		}
	}

	// 将 map 转换为切片
	for _, v := range majorVersionMap {
		versions = append(versions, v)
	}

	// 对版本进行排序（从大到小）
	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i].Version, versions[j].Version) > 0
	})

	return versions, nil
}
