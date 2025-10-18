package cmdCli

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/tea4go/jvms/internal/entity"
	"github.com/tea4go/jvms/utils/jdk"
)

// list 创建列出已安装JDK的命令
// 显示所有已安装的JDK版本，并标记当前正在使用的版本
// 参数:
//   cfx - 配置对象指针
// 返回值:
//   *cli.Command - CLI命令对象指针
func list(cfx *entity.Config) *cli.Command {
	cmd := &cli.Command{
		Name:      "list",
		ShortName: "ls",
		Usage:     "列出当前已安装的JDK",
		Action: func(c *cli.Context) error {
			fmt.Println("已安装的 jdk (* 标记正在使用):")
			v := jdk.GetInstalled(cfx.Store)
			for i, version := range v {
				str := ""
				if cfx.CurrentJDKVersion == version {
					str = fmt.Sprintf("%s  * %d) %s", str, i+1, version)
				} else {
					str = fmt.Sprintf("%s    %d) %s", str, i+1, version)
				}
				fmt.Printf(str + "\n")
			}
			if len(v) == 0 {
				fmt.Println("未识别到已安装的版本。")
			}
			return nil
		},
	}
	return cmd
}
