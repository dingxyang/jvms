package cmdCli

import (
	"github.com/codegangsta/cli"
	"github.com/tea4go/jvms/internal/entity"
)

// use 创建使用JDK版本的命令
// 这是switch命令的别名，功能完全相同
// 参数:
//   cfx - 配置对象指针
// 返回值:
//   *cli.Command - CLI命令对象指针
func use(cfx *entity.Config) *cli.Command {
	cmd := &cli.Command{
		Name:      "use",
		ShortName: "u",
		Usage:     "切换使用指定的版本或索引号",
		Action:    switchFunc(*cfx),
	}
	return cmd
}
