package cmdCli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/tea4go/jvms/internal/entity"
	"github.com/tea4go/jvms/utils/file"
)

// init_ 创建初始化命令
// 用于初始化配置文件和设置环境变量
// 参数:
//   defaultOriginalpath - 默认的JDK下载索引文件URL
//   cfx - 配置对象指针
// 返回值:
//   *cli.Command - CLI命令对象指针
func init_(defaultOriginalpath string, cfx *entity.Config) *cli.Command {
	return &cli.Command{
		Name:        "init",
		Usage:       "初始化配置文件",
		Description: `初始化前请先清空 JAVA_HOME 和 PATH 环境变量。`,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "java_home",
				Usage: "JAVA_HOME 位置",
				Value: filepath.Join(os.Getenv("ProgramFiles"), "jdk"),
			},
			cli.StringFlag{
				Name:  "originalpath",
				Usage: "JDK 下载索引文件 URL",
				Value: defaultOriginalpath,
			},
		},
		Action: func(c *cli.Context) error {
			if c.IsSet("java_home") || cfx.JavaHome == "" {
				cfx.JavaHome = c.String("java_home")
			}
			cmd := exec.Command("cmd", "/C", "setx", "JAVA_HOME", cfx.JavaHome, "/M")
			err := cmd.Run()
			if err != nil {
				return errors.New("设置环境变量 `JAVA_HOME` 失败: 请以管理员身份运行")
			}
			fmt.Println("设置 `JAVA_HOME` 环境变量为 ", cfx.JavaHome)

			if c.IsSet("originalpath") || cfx.Originalpath == "" {
				cfx.Originalpath = c.String("originalpath")
			}
			path := fmt.Sprintf(`%s/bin;%s;%s`, cfx.JavaHome, os.Getenv("PATH"), file.GetCurrentPath())
			cmd = exec.Command("cmd", "/C", "setx", "path", path, "/m")
			err = cmd.Run()
			if err != nil {
				return errors.New("设置环境变量 `PATH` 失败: 请以管理员身份运行")
			}
			fmt.Println("添加 jvms.exe 到 `path` 环境变量")
			return nil
		},
	}

}
