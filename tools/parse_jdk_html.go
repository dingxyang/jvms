package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/tea4go/gh/utils"
	"golang.org/x/net/html"
)

// JDKFile represents a downloadable JDK file
type JDKFile struct {
	Filename     string `json:"filename"`
	URL          string `json:"url"`
	Size         string `json:"size"`
	LastModified string `json:"last_modified"`
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
}

// JDKVersion represents a JDK version entry
type JDKVersion struct {
	Version      string    `json:"version"`
	URL          string    `json:"url"`
	LastModified string    `json:"last_modified"`
	Files        []JDKFile `json:"files"`
}

// JDKData represents the complete output structure
type JDKData struct {
	BaseURL  string       `json:"base_url"`
	Versions []JDKVersion `json:"versions"`
}

func main() {
	// Default URL and output file
	url := "https://mirrors.huaweicloud.com/openjdk/"
	outputFile := "jdk_downloads.json"

	if len(os.Args) > 1 {
		url = os.Args[1]
	}
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}

	fmt.Printf("Fetching JDK versions from: %s\n", url)

	// Fetch HTML content from URL
	var reader io.Reader
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	// Add User-Agent to avoid anti-bot detection
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching URL: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP Error: %s\n", resp.Status)
		os.Exit(1)
	}

	reader = resp.Body

	// Parse HTML
	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		os.Exit(1)
	}

	// Extract JDK versions
	versions := extractJDKVersions(doc)

	// Extract base URL (remove trailing slash)
	baseURL := strings.TrimSuffix(url, "/")

	// Fetch files for each version
	fmt.Printf("\nFetching download files for each version...\n")
	for i := range versions {
		versionURL := baseURL + "/" + versions[i].URL
		fmt.Printf("Fetching files for version %s...\n", versions[i].Version)
		files, err := fetchVersionFiles(versionURL)
		if err != nil {
			fmt.Printf("  Warning: Failed to fetch files for version %s: %v\n", versions[i].Version, err)
			continue
		}
		versions[i].Files = files
		fmt.Printf("  Found %d files\n", len(files))
		time.Sleep(100 * time.Millisecond) // Avoid overwhelming the server
	}

	// Create output data
	jdkData := JDKData{
		BaseURL:  baseURL,
		Versions: versions,
	}
	// Write to JSON file
	jsonData, err := json.MarshalIndent(jdkData, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully parsed %d JDK versions\n", len(versions))
	fmt.Printf("Output written to: %s\n", outputFile)
}

// 从行中提取上次修改日期
func extractLastModified(n *html.Node) string {
	// Move to parent to get the full line text
	if n.Parent != nil && n.Parent.Type == html.ElementNode && n.Parent.Data == "pre" {
		// Get all text from the pre element after the link
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

		// Extract date pattern (e.g., "22-Aug-2021 15:18")
		datePattern := regexp.MustCompile(`(\d{2}-[A-Za-z]{3}-\d{4}\s+\d{2}:\d{2})`)
		matches := datePattern.FindStringSubmatch(text)
		if len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

// 从 HTML 文档中提取 JDK 版本信息
func extractJDKVersions(n *html.Node) []JDKVersion {
	var versions []JDKVersion
	var traverse func(*html.Node)

	// Regex to match version directories (e.g., "10/", "11.0.1/")
	versionPattern := regexp.MustCompile(`^(\d+(\.\d+)*)/$`)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// Get href attribute
			var href string
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href = attr.Val
					break
				}
			}

			// Check if it matches a version pattern
			if versionPattern.MatchString(href) {
				// 使用time.Parse解析日期字符串
				lastModified, err := time.Parse("02-Jan-2006 15:04", extractLastModified(n))
				if err != nil {
					fmt.Println("解析日期错误，", err)
					lastModified = time.Now()
				}
				versions = append(versions,
					JDKVersion{
						Version:      strings.TrimSuffix(href, "/"),
						URL:          href,
						LastModified: lastModified.Format(utils.DateTimeFormat),
					})
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)
	return versions
}

// fetchVersionFiles fetches the list of downloadable files for a specific JDK version
func fetchVersionFiles(versionURL string) ([]JDKFile, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", versionURL, nil)
	if err != nil {
		return nil, err
	}

	// Add User-Agent to avoid anti-bot detection
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %s", resp.Status)
	}

	// Parse HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract files
	return extractFiles(doc, versionURL), nil
}

// extractFiles extracts downloadable files from version directory HTML
func extractFiles(n *html.Node, baseURL string) []JDKFile {
	var files []JDKFile
	var traverse func(*html.Node)

	// Regex to match file extensions (tar.gz, zip, etc.)
	filePattern := regexp.MustCompile(`\.(tar\.gz|zip|msi|pkg|bin)$`)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// Get href attribute
			var href string
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href = attr.Val
					break
				}
			}

			// Check if it's a downloadable file
			if filePattern.MatchString(href) {
				// Extract file info
				lastModified := extractLastModified(n)
				size := extractSize(n)

				// Parse date
				parsedDate, err := time.Parse("02-Jan-2006 15:04", lastModified)
				if err != nil {
					parsedDate = time.Now()
				}

				// Parse GOOS and GOARCH from filename
				goos, goarch := parseOSAndArch(href)

				files = append(files, JDKFile{
					Filename:     href,
					URL:          baseURL + "/" + href,
					Size:         size,
					LastModified: parsedDate.Format(utils.DateTimeFormat),
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

// extractSize extracts file size from the HTML node context
func extractSize(n *html.Node) string {
	// Move to parent to get the full line text
	if n.Parent != nil && n.Parent.Type == html.ElementNode && n.Parent.Data == "pre" {
		// Get all text from the pre element after the link
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

		// Extract size pattern (e.g., "123.45M" or "1.2G")
		sizePattern := regexp.MustCompile(`\s+([\d.]+[KMG]?)\s*$`)
		matches := sizePattern.FindStringSubmatch(text)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}
	return ""
}

// parseOSAndArch extracts GOOS and GOARCH from JDK filename
// Examples:
//   - openjdk-24_linux-x64_bin.tar.gz -> (linux, amd64)
//   - openjdk-24_windows-x64_bin.zip -> (windows, amd64)
//   - openjdk-24_macos-aarch64_bin.tar.gz -> (darwin, arm64)
func parseOSAndArch(filename string) (goos, goarch string) {
	filename = strings.ToLower(filename)

	// Mapping from JDK OS names to GOOS
	osMap := map[string]string{
		"linux":   "linux",
		"windows": "windows",
		"macos":   "darwin",
		"osx":     "darwin",
		"darwin":  "darwin",
		"solaris": "solaris",
		"aix":     "aix",
	}

	// Mapping from JDK architecture names to GOARCH
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

	// Extract OS
	for jdkOS, goosName := range osMap {
		if strings.Contains(filename, jdkOS) {
			goos = goosName
			break
		}
	}

	// Extract architecture
	for jdkArch, goarchName := range archMap {
		if strings.Contains(filename, jdkArch) {
			goarch = goarchName
			break
		}
	}

	// If not found, set to unknown
	if goos == "" {
		goos = "unknown"
	}
	if goarch == "" {
		goarch = "unknown"
	}

	return goos, goarch
}
