package cmdCli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/tea4go/jvms/internal/entity"
	"github.com/tea4go/jvms/utils/web"
)

// rlsCmd 执行显示可下载版本列表的命令
// 显示可供下载的JDK版本列表
// 参数:
//
//	args - 命令参数
//	cfx - 配置对象指针
//
// 返回值:
//
//	error - 执行错误
func rlsCmd(args []string, cfx *entity.TConfig) error {
	// 创建命令专用的 FlagSet
	fs := pflag.NewFlagSet("rls", pflag.ContinueOnError)

	showAll := fs.BoolP("a", "a", false, "列出所有版本")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if cfx.Proxy != "" {
		web.SetProxy(cfx.Proxy)
	}

	versions, err := getJdkVersions(cfx)
	if err != nil {
		return err
	}

	for i, version := range versions {
		fmt.Printf("    %d) %s\n", i+1, version.Version)
		if !*showAll && i >= 999 {
			fmt.Println("\n使用 \"jvm rls -a\" 显示所有版本")
			break
		}
	}

	if len(versions) == 0 {
		fmt.Println("没有可供下载的 jdk 版本。")
	}

	fmt.Printf("\n完整列表请访问 %s\n", cfx.Originalpath)
	return nil
}
