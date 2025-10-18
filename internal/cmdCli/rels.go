package cmdCli

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/tea4go/jvms/internal/entity"
	"github.com/tea4go/jvms/utils/web"
)

// rls 创建显示可下载版本列表的命令
// 显示可供下载的JDK版本列表
// 参数:
//   cfx - 配置对象指针
// 返回值:
//   *cli.Command - CLI命令对象指针
func rls(cfx *entity.Config) *cli.Command {
	cmd := &cli.Command{
		Name:  "rls",
		Usage: "显示可供下载的版本列表",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "a",
				Usage: "列出所有版本",
			},
		},
		Action: func(c *cli.Context) error {
			if cfx.Proxy != "" {
				web.SetProxy(cfx.Proxy)
			}
			versions, err := getJdkVersions(cfx)
			if err != nil {
				return err
			}
			for i, version := range versions {
				fmt.Printf("    %d) %s\n", i+1, version.Version)
				if !c.Bool("a") && i >= 9 {
					fmt.Println("\n使用 \"jvm rls -a\" 显示所有版本")
					break
				}
			}
			if len(versions) == 0 {
				fmt.Println("没有可供下载的 jdk 版本。")
			}

			fmt.Printf("\n完整列表请访问 %s\n", cfx.Originalpath)
			return nil
		},
	}
	return cmd
}
