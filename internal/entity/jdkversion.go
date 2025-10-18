package entity

// JdkVersion JDK版本信息结构体，用于存储JDK的版本号和下载地址
type JdkVersion struct {
	// Version JDK版本号
	Version string `json:"version"`
	// Url JDK下载地址
	Url string `json:"url"`
}
