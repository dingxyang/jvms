package jdk

import (
	"fmt"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// HuaweiJDK 表示华为镜像源中的 JDK 文件信息
type HuaweiJDK struct {
	Version      string // 版本号
	Filename     string // 文件名
	URL          string // 下载链接
	Size         string // 文件大小
	LastModified string // 最后修改时间
	GOOS         string // 操作系统类型
	GOARCH       string // 系统架构
}

// HuaweiJDKs 从华为镜像源获取所有 JDK 版本
// 返回值:
//   []HuaweiJDK - JDK 版本信息列表
func HuaweiJDKs() []HuaweiJDK {
	return HuaweiJDKsFromURL("https://mirrors.huaweicloud.com/openjdk/")
}

// HuaweiJDKsFromURL 从指定 URL 获取 JDK 版本列表
// 参数:
//   baseURL - 镜像源的基础 URL
// 返回值:
//   []HuaweiJDK - JDK 版本信息列表
func HuaweiJDKsFromURL(baseURL string) []HuaweiJDK {
	var allJDKs []HuaweiJDK

	// 获取版本列表
	versions := fetchVersionList(baseURL)

	// 获取每个版本的文件列表
	for _, version := range versions {
		versionURL := baseURL + version
		files := fetchVersionFiles(versionURL, version)

		// 过滤出与当前运行时操作系统和架构匹配的文件
		for _, file := range files {
			if file.GOOS == runtime.GOOS && file.GOARCH == runtime.GOARCH {
				allJDKs = append(allJDKs, file)
			}
		}

		time.Sleep(50 * time.Millisecond) // 避免对服务器造成过大压力
	}

	return allJDKs
}

// fetchVersionList 获取 JDK 版本目录列表
// 参数:
//   baseURL - 镜像源的基础 URL
// 返回值:
//   []string - 版本目录名称列表
func fetchVersionList(baseURL string) []string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		fmt.Printf("创建请求时出错: %v\n", err)
		return nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("获取 URL 时出错: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP 错误: %s\n", resp.Status)
		return nil
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Printf("解析 HTML 时出错: %v\n", err)
		return nil
	}

	return extractVersions(doc)
}

// extractVersions 从 HTML 中提取版本目录
// 参数:
//   n - HTML 节点
// 返回值:
//   []string - 版本目录名称列表
func extractVersions(n *html.Node) []string {
	var versions []string
	versionPattern := regexp.MustCompile(`^(\d+(\.\d+)*)/$`)

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var href string
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href = attr.Val
					break
				}
			}

			if versionPattern.MatchString(href) {
				versions = append(versions, href)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return versions
}

// fetchVersionFiles 获取特定版本的文件列表
// 参数:
//   versionURL - 版本目录的 URL
//   version - 版本号
// 返回值:
//   []HuaweiJDK - 该版本的 JDK 文件信息列表
func fetchVersionFiles(versionURL string, version string) []HuaweiJDK {
	client := &http.Client{}
	req, err := http.NewRequest("GET", versionURL, nil)
	if err != nil {
		return nil
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil
	}

	return extractFiles(doc, versionURL, version)
}

// extractFiles 从版本目录的 HTML 中提取可下载的文件
// 参数:
//   n - HTML 节点
//   baseURL - 基础 URL
//   version - 版本号
// 返回值:
//   []HuaweiJDK - JDK 文件信息列表
func extractFiles(n *html.Node, baseURL string, version string) []HuaweiJDK {
	var files []HuaweiJDK
	filePattern := regexp.MustCompile(`\.(tar\.gz|zip|msi|pkg|bin)$`)

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var href string
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href = attr.Val
					break
				}
			}

			if filePattern.MatchString(href) {
				lastModified := extractLastModified(n)
				size := extractSize(n)
				goos, goarch := parseOSAndArch(href)

				files = append(files, HuaweiJDK{
					Version:      strings.TrimSuffix(version, "/"),
					Filename:     href,
					URL:          baseURL + href,
					Size:         size,
					LastModified: lastModified,
					GOOS:         goos,
					GOARCH:       goarch,
				})
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return files
}

// extractLastModified 从 HTML 节点上下文中提取最后修改日期
// 参数:
//   n - HTML 节点
// 返回值:
//   string - 最后修改日期字符串，未找到则返回空字符串
func extractLastModified(n *html.Node) string {
	if n.Parent != nil && n.Parent.Type == html.ElementNode && n.Parent.Data == "pre" {
		var found bool
		var text string
		for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
			if c == n {
				found = true
				continue
			}
			if found && c.Type == html.TextNode {
				text += c.Data
				break
			}
		}

		datePattern := regexp.MustCompile(`(\d{2}-[A-Za-z]{3}-\d{4}\s+\d{2}:\d{2})`)
		matches := datePattern.FindStringSubmatch(text)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

// extractSize 从 HTML 节点上下文中提取文件大小
// 参数:
//   n - HTML 节点
// 返回值:
//   string - 文件大小字符串，未找到则返回空字符串
func extractSize(n *html.Node) string {
	if n.Parent != nil && n.Parent.Type == html.ElementNode && n.Parent.Data == "pre" {
		var found bool
		var text string
		for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
			if c == n {
				found = true
				continue
			}
			if found && c.Type == html.TextNode {
				text += c.Data
				break
			}
		}

		sizePattern := regexp.MustCompile(`\s+([\d.]+[KMG]?)\s*$`)
		matches := sizePattern.FindStringSubmatch(text)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}
	return ""
}

// parseOSAndArch 从 JDK 文件名中提取操作系统和架构信息
// 参数:
//   filename - JDK 文件名
// 返回值:
//   goos - 操作系统类型（GOOS 格式）
//   goarch - 系统架构（GOARCH 格式）
func parseOSAndArch(filename string) (goos, goarch string) {
	filename = strings.ToLower(filename)

	// JDK 操作系统名称到 GOOS 的映射
	osMap := map[string]string{
		"linux":   "linux",
		"windows": "windows",
		"macos":   "darwin",
		"osx":     "darwin",
		"darwin":  "darwin",
		"solaris": "solaris",
		"aix":     "aix",
	}

	// JDK 架构名称到 GOARCH 的映射
	archMap := map[string]string{
		"x64":     "amd64",
		"x86_64":  "amd64",
		"amd64":   "amd64",
		"aarch64": "arm64",
		"arm64":   "arm64",
		"x86":     "386",
		"i386":    "386",
		"i686":    "386",
		"ppc64":   "ppc64",
		"ppc64le": "ppc64le",
		"s390x":   "s390x",
	}

	// 提取操作系统
	for jdkOS, goosName := range osMap {
		if strings.Contains(filename, jdkOS) {
			goos = goosName
			break
		}
	}

	// 提取架构
	for jdkArch, goarchName := range archMap {
		if strings.Contains(filename, jdkArch) {
			goarch = goarchName
			break
		}
	}

	// 如果未找到，设置为 unknown
	if goos == "" {
		goos = "unknown"
	}
	if goarch == "" {
		goarch = "unknown"
	}

	return goos, goarch
}
