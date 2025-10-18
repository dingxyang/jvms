package cmdCli

import (
	"fmt"

	"github.com/tea4go/jvms/internal/entity"
)

// TCommandParams 命令参数结构体
// 包含默认原始路径和配置对象
type TCommandParams struct {
	Config *entity.TConfig // 配置对象指针
}

// Execute 执行指定的命令
// 参数:
//
//	command - 命令名称
//	args - 命令参数
//	cp - 命令参数对象指针
//
// 返回值:
//
//	error - 执行错误
func Execute(command string, args []string, cp *TCommandParams) error {
	switch command {
	case "init":
		return initCmd(args, cp.Config)
	case "list", "ls":
		return listCmd(args, cp.Config)
	case "install", "i":
		return installCmd(args, cp.Config)
	case "switch", "s":
		return switchCmd(args, cp.Config)
	case "use", "u":
		return useCmd(args, cp.Config)
	case "remove", "rm":
		return removeCmd(args, cp.Config)
	case "rls":
		return rlsCmd(args, cp.Config)
	case "proxy":
		return proxyCmd(args, cp.Config)
	case "help", "h":
		return printHelp(args)
	default:
		return fmt.Errorf("未知命令: %s\n使用 'jvms help' 查看可用命令", command)
	}
}

// printHelp 打印帮助信息
func printHelp(args []string) error {
	if len(args) == 0 {
		fmt.Println("JVMS - JDK Version Manager for Windows")
		fmt.Println("")
		fmt.Println("可用命令:")
		fmt.Println("  init        初始化配置文件")
		fmt.Println("  list, ls    列出当前已安装的JDK")
		fmt.Println("  install, i  安装可用的远程JDK")
		fmt.Println("  switch, s   切换使用指定的版本或索引号")
		fmt.Println("  use, u      切换使用指定的版本或索引号")
		fmt.Println("  remove, rm  删除指定的版本")
		fmt.Println("  rls         显示可供下载的版本列表")
		fmt.Println("  proxy       设置下载使用的代理")
		fmt.Println("")
		fmt.Println("使用 'jvms help <命令>' 查看命令的详细帮助")
		return nil
	}

	// 这里可以添加针对具体命令的详细帮助
	cmd := args[0]
	switch cmd {
	case "init":
		fmt.Println("init - 初始化配置文件")
		fmt.Println("")
		fmt.Println("用法: jvms init [选项]")
		fmt.Println("")
		fmt.Println("选项:")
		fmt.Println("  --java_home <路径>      指定 JAVA_HOME 位置")
	case "install", "i":
		fmt.Println("install - 安装可用的远程JDK")
		fmt.Println("")
		fmt.Println("用法: jvms install <版本>")
	case "list", "ls":
		fmt.Println("list - 列出当前已安装的JDK")
		fmt.Println("")
		fmt.Println("用法: jvms list")
	case "switch", "s":
		fmt.Println("switch - 切换使用指定的版本或索引号")
		fmt.Println("")
		fmt.Println("用法: jvms switch <版本|索引号>")
	case "use", "u":
		fmt.Println("use - 切换使用指定的版本或索引号")
		fmt.Println("")
		fmt.Println("用法: jvms use <版本|索引号>")
	case "remove", "rm":
		fmt.Println("remove - 删除指定的版本")
		fmt.Println("")
		fmt.Println("用法: jvms remove <版本>")
	case "rls":
		fmt.Println("rls - 显示可供下载的版本列表")
		fmt.Println("")
		fmt.Println("用法: jvms rls [-a]")
		fmt.Println("")
		fmt.Println("选项:")
		fmt.Println("  -a    列出所有版本")
	case "proxy":
		fmt.Println("proxy - 设置下载使用的代理")
		fmt.Println("")
		fmt.Println("用法: jvms proxy [选项]")
		fmt.Println("")
		fmt.Println("选项:")
		fmt.Println("  --show          显示当前代理")
		fmt.Println("  --set <代理>    设置代理")
	default:
		fmt.Printf("未知命令: %s\n", cmd)
	}
	return nil
}
