package cmdCli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/tea4go/jvms/internal/entity"
)

// proxyCmd 执行代理设置命令
// 设置或显示用于下载的代理服务器
// 参数:
//   args - 命令参数
//   cfx - 配置对象指针
// 返回值:
//   error - 执行错误
func proxyCmd(args []string, cfx *entity.TConfig) error {
	// 创建命令专用的 FlagSet
	fs := pflag.NewFlagSet("proxy", pflag.ContinueOnError)

	show := fs.Bool("show", false, "显示代理")
	set := fs.String("set", "", "设置代理")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *show {
		fmt.Printf("当前代理: %s\n", cfx.Proxy)
		return nil
	}

	if fs.Changed("set") {
		cfx.Proxy = *set
		fmt.Printf("代理已设置为: %s\n", cfx.Proxy)
	}

	return nil
}
