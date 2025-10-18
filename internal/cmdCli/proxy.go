package cmdCli

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/tea4go/jvms/internal/entity"
)

// proxy 创建代理设置命令
// 设置或显示用于下载的代理服务器
// 参数:
//   cfx - 配置对象指针
// 返回值:
//   *cli.Command - CLI命令对象指针
func proxy(cfx *entity.TConfig) *cli.Command {
	cmd := &cli.Command{
		Name:  "proxy",
		Usage: "设置下载使用的代理",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "show",
				Usage: "显示代理",
			},
			cli.StringFlag{
				Name:  "set",
				Usage: "设置代理",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("show") {
				fmt.Printf("当前代理: %s\n", cfx.Proxy)
				return nil
			}
			if c.IsSet("set") {
				cfx.Proxy = c.String("set")
			}
			return nil
		},
	}
	return cmd
}
