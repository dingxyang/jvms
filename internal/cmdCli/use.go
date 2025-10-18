package cmdCli

import (
	"github.com/tea4go/jvms/internal/entity"
)

// useCmd 执行使用JDK版本的命令
// 这是switch命令的别名，功能完全相同
// 参数:
//   args - 命令参数
//   cfx - 配置对象指针
// 返回值:
//   error - 执行错误
func useCmd(args []string, cfx *entity.TConfig) error {
	return switchFunc(args, cfx)
}
