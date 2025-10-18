package cmdCli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/tea4go/jvms/internal/entity"
	"github.com/tea4go/jvms/utils/jdk"
)

// remove 创建删除JDK的命令
// 删除指定版本的JDK安装
// 参数:
//   cfx - 配置对象指针
// 返回值:
//   *cli.Command - CLI命令对象指针
func remove(cfx *entity.TConfig) *cli.Command {
	cmd := &cli.Command{
		Name:      "remove",
		ShortName: "rm",
		Usage:     "删除指定的版本",
		Action: func(c *cli.Context) error {
			v := c.Args().Get(0)
			if v == "" {
				return errors.New("您应该输入一个版本，输入 \"jvms list\" 查看已安装的版本")
			}
			if jdk.IsVersionInstalled(cfx.Store, v) {
				fmt.Printf("删除 JDK %s ...\n", v)
				if cfx.CurrentJDKVersion == v {
					os.Remove(cfx.JavaHome)
				}
				dir := filepath.Join(cfx.Store, v)
				e := os.RemoveAll(dir)
				if e != nil {
					fmt.Println("删除 jdk " + v + " 时出错")
					fmt.Println("请手动删除 " + dir + "。")
				} else {
					fmt.Printf(" 完成")
				}
			} else {
				fmt.Println("jdk " + v + " 未安装。输入 \"jvms list\" 查看已安装的版本。")
			}
			return nil
		},
	}
	return cmd
}
