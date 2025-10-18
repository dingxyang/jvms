# JDK 下载地址解析器

## 概述

`parse_jdk_html.go` 是一个自动化工具，用于从华为云镜像站点（或其他OpenJDK镜像站）抓取所有JDK版本及其下载文件的完整信息，并输出为结构化的JSON文件。

## 功能特性

### 核心功能

1. **自动获取版本列表** - 直接从URL获取所有可用的JDK版本
2. **递归解析下载文件** - 进入每个版本目录，提取真实的下载文件信息
3. **完整元数据采集** - 收集文件名、完整URL、文件大小、最后修改时间
4. **JSON结构化输出** - 生成格式化的JSON文件，便于程序化使用

### 技术特点

- **HTTP请求优化** - 添加User-Agent模拟浏览器，避免反爬虫机制
- **请求节流** - 每个请求间隔100ms，避免对服务器造成压力
- **错误处理** - 完善的错误处理和友好的进度提示
- **灵活配置** - 支持命令行参数自定义URL和输出文件

## 使用方法

### 基本用法

```bash
# 使用默认配置（华为云镜像）
go run parse_jdk_html.go

# 自定义URL
go run parse_jdk_html.go https://example.com/openjdk/

# 自定义URL和输出文件
go run parse_jdk_html.go https://mirrors.huaweicloud.com/openjdk/ custom_output.json
```

### 编译运行

```bash
# 编译
go build parse_jdk_html.go

# 运行
./parse_jdk_html

# Windows
parse_jdk_html.exe
```

## 输出格式

### JSON结构

```json
{
  "base_url": "https://mirrors.huaweicloud.com/openjdk",
  "versions": [
    {
      "version": "24",
      "url": "24/",
      "last_modified": "2025-02-07 02:46:00",
      "files": [
        {
          "filename": "openjdk-24_linux-x64_bin.tar.gz",
          "url": "https://mirrors.huaweicloud.com/openjdk/24//openjdk-24_linux-x64_bin.tar.gz",
          "size": "",
          "last_modified": "2025-02-07 02:46:00"
        },
        {
          "filename": "openjdk-24_windows-x64_bin.zip",
          "url": "https://mirrors.huaweicloud.com/openjdk/24//openjdk-24_windows-x64_bin.zip",
          "size": "",
          "last_modified": "2025-02-07 02:46:00",
          "goos": "windows",
          "goarch": "amd64"
        }
      ]
    }
  ]
}
```

### 数据字段说明

#### 根对象
- `base_url`: 镜像站点的基础URL
- `versions`: JDK版本数组

#### Version对象
- `version`: 版本号（如 "24", "17.0.1"）
- `url`: 版本目录的相对路径
- `last_modified`: 版本目录的最后修改时间
- `files`: 该版本下所有下载文件的数组

#### File对象
- `filename`: 文件名（如 "openjdk-24_linux-x64_bin.tar.gz"）
- `url`: 文件的完整下载URL
- `size`: 文件大小（可能为空）
- `last_modified`: 文件的最后修改时间
- `goos`: Go语言的操作系统标识（linux, windows, darwin等）
- `goarch`: Go语言的架构标识（amd64, arm64, 386等）

## 解析统计

基于华为云镜像站点的解析结果：

- **JDK版本数**: 50个（JDK 9 - JDK 25）
- **下载文件数**: 204个
- **支持平台**:
  - Linux (x64, aarch64) - 83个文件
  - macOS/Darwin (x64, aarch64) - 73个文件
  - Windows (x64) - 48个文件
- **架构分布**:
  - amd64/x64: 145个文件
  - arm64/aarch64: 59个文件
- **文件格式**: tar.gz, zip, msi, pkg, bin

## 依赖项

```go
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
```

### 安装依赖

```bash
go get golang.org/x/net/html
```

## 代码结构

### 主要数据类型

```go
// JDKFile 表示可下载的JDK文件
type JDKFile struct {
    Filename     string `json:"filename"`
    URL          string `json:"url"`
    Size         string `json:"size"`
    LastModified string `json:"last_modified"`
    GOOS         string `json:"goos"`
    GOARCH       string `json:"goarch"`
}

// JDKVersion 表示JDK版本条目
type JDKVersion struct {
    Version      string    `json:"version"`
    URL          string    `json:"url"`
    LastModified string    `json:"last_modified"`
    Files        []JDKFile `json:"files"`
}

// JDKData 表示完整的输出结构
type JDKData struct {
    BaseURL  string       `json:"base_url"`
    Versions []JDKVersion `json:"versions"`
}
```

### 主要函数

- `main()` - 程序入口，协调整体流程
- `extractJDKVersions(n *html.Node)` - 从主页面提取版本列表
- `fetchVersionFiles(versionURL string)` - 获取指定版本的文件列表
- `extractFiles(n *html.Node, baseURL string)` - 从版本页面提取文件信息
- `extractLastModified(n *html.Node)` - 提取最后修改时间
- `extractSize(n *html.Node)` - 提取文件大小
- `parseOSAndArch(filename string)` - 从文件名解析GOOS和GOARCH

## 工作流程

1. **初始化** - 设置默认URL和输出文件名，解析命令行参数
2. **获取版本列表** - 发送HTTP请求到主URL，解析HTML获取所有版本目录
3. **遍历版本** - 对每个版本：
   - 构造版本目录的完整URL
   - 发送HTTP请求获取该版本的HTML
   - 解析HTML提取所有下载文件信息
   - 添加100ms延迟
4. **生成JSON** - 将所有数据序列化为格式化的JSON
5. **写入文件** - 保存到指定的输出文件

## 正则表达式模式

- **版本目录匹配**: `^(\d+(\.\d+)*)/$` - 匹配如 "24/", "17.0.1/" 的目录
- **文件扩展名匹配**: `\.(tar\.gz|zip|msi|pkg|bin)$` - 匹配可下载文件
- **日期模式**: `(\d{2}-[A-Za-z]{3}-\d{4}\s+\d{2}:\d{2})` - 如 "07-Feb-2025 02:46"
- **大小模式**: `\s+([\d.]+[KMG]?)\s*$` - 如 "123.45M"

## GOOS和GOARCH解析

程序会自动从文件名中提取操作系统和架构信息，并映射到Go语言的标准标识符。

### 操作系统映射 (GOOS)

| JDK文件名标识 | GOOS值 | 说明 |
|-------------|--------|------|
| linux | linux | Linux系统 |
| windows | windows | Windows系统 |
| macos | darwin | macOS系统（新版本） |
| osx | darwin | macOS系统（旧版本） |
| darwin | darwin | macOS系统 |
| solaris | solaris | Solaris系统 |
| aix | aix | AIX系统 |

### 架构映射 (GOARCH)

| JDK文件名标识 | GOARCH值 | 说明 |
|-------------|----------|------|
| x64, x86_64, amd64 | amd64 | 64位x86架构 |
| aarch64, arm64 | arm64 | 64位ARM架构 |
| x86, i386, i686 | 386 | 32位x86架构 |
| ppc64 | ppc64 | 64位PowerPC架构 |
| ppc64le | ppc64le | 64位PowerPC小端架构 |
| s390x | s390x | IBM z架构 |

### 解析示例

```
openjdk-24_linux-x64_bin.tar.gz
  → goos: linux, goarch: amd64

openjdk-24_windows-x64_bin.zip
  → goos: windows, goarch: amd64

openjdk-24_macos-aarch64_bin.tar.gz
  → goos: darwin, goarch: arm64

openjdk-10_osx-x64_bin.tar.gz
  → goos: darwin, goarch: amd64
```

## 注意事项

1. **网络依赖** - 需要稳定的网络连接访问镜像站点
2. **执行时间** - 由于需要访问50+个页面，完整执行需要约5-10秒
3. **User-Agent** - 使用浏览器User-Agent避免被反爬虫机制拦截
4. **请求频率** - 已添加100ms延迟，避免触发服务器限流

## 扩展用途

此工具生成的JSON文件可用于：

- **自动化下载** - 根据平台和版本自动选择并下载JDK
  ```go
  // 示例：根据运行时环境选择合适的JDK
  goos := runtime.GOOS
  goarch := runtime.GOARCH
  // 从JSON中查找匹配的文件
  ```
- **版本管理** - 构建JDK版本管理工具
- **镜像监控** - 监控镜像站点的更新情况
- **统计分析** - 分析JDK版本发布历史和趋势
- **跨平台部署** - 基于GOOS/GOARCH自动选择适合目标平台的JDK包

## 示例输出

运行程序时的控制台输出：

```
Fetching JDK versions from: https://mirrors.huaweicloud.com/openjdk/

Fetching download files for each version...
Fetching files for version 10...
  Found 3 files
Fetching files for version 10.0.1...
  Found 3 files
...
Fetching files for version 25...
  Found 5 files
Successfully parsed 50 JDK versions
Output written to: jdk_downloads.json
```

## 测试和验证

### 查看生成的JSON文件

```bash
# 查看文件前80行
head -80 jdk_downloads.json

# 查看特定版本的信息
cat jdk_downloads.json | grep -A 40 '"version": "24"'

# 使用jq查看格式化输出
cat jdk_downloads.json | jq '.'
```

### 统计信息

```bash
# 统计版本数
cat jdk_downloads.json | grep -c '"version":'

# 统计文件数
cat jdk_downloads.json | grep -c '"filename":'

# 查看GOOS分布
cat jdk_downloads.json | grep '"goos":' | sort | uniq -c

# 查看GOARCH分布
cat jdk_downloads.json | grep '"goarch":' | sort | uniq -c
```

### 验证GOOS和GOARCH解析

```bash
# 检查是否有未识别的GOOS
cat jdk_downloads.json | grep '"goos": "unknown"' | wc -l

# 检查是否有未识别的GOARCH
cat jdk_downloads.json | grep '"goarch": "unknown"' | wc -l

# 查看特定版本的所有文件及其平台信息（需要安装jq）
cat jdk_downloads.json | jq '.versions[] | select(.version == "17") | .files[] | {filename, goos, goarch}'

# 查看Windows平台的文件
cat jdk_downloads.json | grep -A 7 'windows'

# 查看Linux平台的文件
cat jdk_downloads.json | grep -A 7 'linux'

# 查看macOS平台的文件
cat jdk_downloads.json | grep -A 7 'darwin'
```

### 查询特定平台的JDK

```bash
# 查找适合当前系统的JDK 17文件（使用jq）
cat jdk_downloads.json | jq --arg goos "$(go env GOOS)" --arg goarch "$(go env GOARCH)" \
  '.versions[] | select(.version == "17") | .files[] | select(.goos == $goos and .goarch == $goarch)'

# 查找所有Linux amd64平台的文件
cat jdk_downloads.json | jq '.versions[].files[] | select(.goos == "linux" and .goarch == "amd64") | .filename'

# 查找所有ARM64架构的文件
cat jdk_downloads.json | jq '.versions[].files[] | select(.goarch == "arm64") | {version: .filename, url: .url}'
```

### 统计分析

```bash
# 统计每个版本的文件数
cat jdk_downloads.json | jq '.versions[] | {version: .version, file_count: (.files | length)}'

# 统计每个操作系统的文件总数
echo "=== GOOS分布 ==="
cat jdk_downloads.json | grep '"goos":' | sort | uniq -c

echo ""
echo "=== GOARCH分布 ==="
cat jdk_downloads.json | grep '"goarch":' | sort | uniq -c

# 查看最新版本的信息
cat jdk_downloads.json | jq '.versions[] | select(.version == "25")'
```

### 测试结果示例

```bash
# GOOS分布结果
=== GOOS分布 ===
     73           "goos": "darwin",
     83           "goos": "linux",
     48           "goos": "windows",

# GOARCH分布结果
=== GOARCH分布 ===
    145           "goarch": "amd64"
     59           "goarch": "arm64"

# 验证结果
未识别的GOOS数量: 0
未识别的GOARCH数量: 0
```

## 许可证

本程序基于项目许可证开源。

## 维护

如遇到问题或需要改进，请提交Issue或Pull Request。
