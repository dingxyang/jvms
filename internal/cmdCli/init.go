package cmdCli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/tea4go/jvms/internal/entity"
	"github.com/tea4go/jvms/utils/file"
)

// initCmd 执行初始化命令
// 用于初始化配置文件和设置环境变量
// 参数:
//   args - 命令参数
//   defaultOriginalpath - 默认的JDK下载索引文件URL
//   cfx - 配置对象指针
// 返回值:
//   error - 执行错误
func initCmd(args []string, defaultOriginalpath string, cfx *entity.TConfig) error {
	// 创建命令专用的 FlagSet
	fs := pflag.NewFlagSet("init", pflag.ContinueOnError)

	javaHome := fs.String("java_home", filepath.Join(os.Getenv("ProgramFiles"), "jdk"), "JAVA_HOME 位置")
	originalpath := fs.String("originalpath", defaultOriginalpath, "JDK 下载索引文件 URL")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// 设置 JAVA_HOME
	if fs.Changed("java_home") || cfx.JavaHome == "" {
		cfx.JavaHome = *javaHome
	}

	cmd := exec.Command("cmd", "/C", "setx", "JAVA_HOME", cfx.JavaHome, "/M")
	err := cmd.Run()
	if err != nil {
		return errors.New("设置环境变量 `JAVA_HOME` 失败: 请以管理员身份运行")
	}
	fmt.Println("设置 `JAVA_HOME` 环境变量为 ", cfx.JavaHome)

	// 设置 originalpath
	if fs.Changed("originalpath") || cfx.Originalpath == "" {
		cfx.Originalpath = *originalpath
	}

	// 设置 PATH
	path := fmt.Sprintf(`%s/bin;%s;%s`, cfx.JavaHome, os.Getenv("PATH"), file.GetCurrentPath())
	cmd = exec.Command("cmd", "/C", "setx", "path", path, "/m")
	err = cmd.Run()
	if err != nil {
		return errors.New("设置环境变量 `PATH` 失败: 请以管理员身份运行")
	}
	fmt.Println("添加 jvms.exe 到 `path` 环境变量")

	return nil
}
