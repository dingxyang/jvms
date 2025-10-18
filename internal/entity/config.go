package entity

// TConfig 配置结构体，用于存储JVMS的全局配置信息
type TConfig struct {
	// JavaHome Java环境变量路径
	JavaHome string `json:"java_home"`
	// CurrentJDKVersion 当前使用的JDK版本
	CurrentJDKVersion string `json:"current_jdk_version"`
	// Originalpath 原始的PATH环境变量值
	Originalpath string `json:"original_path"`
	// Proxy 代理服务器地址
	Proxy string `json:"proxy"`
	// Store JDK存储路径
	Store string
	// Download JDK下载路径
	Download string
}
