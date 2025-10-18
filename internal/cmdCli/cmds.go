package cmdCli

import (
	"github.com/codegangsta/cli"
	"github.com/tea4go/jvms/internal/entity"
)

// CommandParams 命令参数结构体
// 包含默认原始路径和配置对象
type CommandParams struct {
	DefaultOriginalPath string         // 默认原始路径
	Config              *entity.Config // 配置对象指针
}

// Commands 获取所有可用的CLI命令列表
// 参数:
//   cp - 命令参数对象指针
// 返回值:
//   []cli.Command - CLI命令列表
func Commands(cp *CommandParams) []cli.Command {
	cmds := []cli.Command{
		*init_(cp.DefaultOriginalPath, cp.Config),
		*list(cp.Config),
		*install(cp.Config),
		*switch_(cp.Config),
		*use(cp.Config),
		*remove(cp.Config),
		*rls(cp.Config),
		*proxy(cp.Config),
	}
	return cmds
}
