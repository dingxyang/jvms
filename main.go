// Package main 是 JVMS (JDK Version Manager for Windows) 的主入口
// JVMS 是一个用于在 Windows 系统上管理多个 JDK 版本的命令行工具
package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/tea4go/jvms/internal/cmdCli"
	"github.com/tea4go/jvms/internal/entity"
	"github.com/tea4go/jvms/utils/file"
	"github.com/tea4go/jvms/utils/web"
	"github.com/tucnak/store"
)

// version 定义当前 JVMS 的版本号
var version = "2.1.0"

const (
	// defaultOriginalpath 定义默认的 JDK 下载索引文件 URL
	// 注意：此 URL 已不再使用，实际使用华为云镜像
	defaultOriginalpath = "https://raw.githubusercontent.com/ystyle/jvms/new/jdkdlindex.json"
)

// cfx 全局配置对象，存储 JVMS 的运行配置
var cfx entity.TConfig

// main 是程序的入口函数
// 初始化 CLI 应用并执行用户命令
func main() {
	app := cli.NewApp()
	app.Name = "jvms"
	app.Usage = `JDK Version Manager (JVMS) for Windows`
	app.Version = version
	app.CommandNotFound = commandNotFound
	app.Commands = commands()

	// 设置启动前和关闭后的钩子函数
	app.Before = startup
	app.After = shutdown

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

// commands 返回所有可用的 CLI 命令列表
// 包括 init, list, install, switch, use, remove, rls, proxy 等命令
func commands() []cli.Command {
	cmds := cmdCli.Commands(&cmdCli.TCommandParams{
		DefaultOriginalPath: defaultOriginalpath,
		Config:              &cfx,
	})
	return cmds
}

// commandNotFound 处理用户输入了不存在的命令的情况
// 参数:
//   c - CLI 上下文
//   command - 用户输入的命令名称
func commandNotFound(c *cli.Context, command string) {
	log.Fatal("Command Not Found")
}

// startup 在应用启动前执行
// 主要功能：
//   1. 注册 JSON 序列化/反序列化器
//   2. 加载配置文件 (jvms.json)
//   3. 初始化存储路径和下载路径
//   4. 设置代理（如果配置了）
// 参数:
//   c - CLI 上下文
// 返回:
//   error - 初始化失败时返回错误
func startup(c *cli.Context) error {
	// 注册 JSON 格式的配置存储器
	store.Register(
		"json",
		func(v interface{}) ([]byte, error) {
			return json.MarshalIndent(v, "", "    ")
		},
		json.Unmarshal)

	// 初始化配置存储
	store.Init("jvms")

	// 加载配置文件
	if err := store.Load("jvms.json", &cfx); err != nil {
		return errors.New("failed to load the config:" + err.Error())
	}

	// 获取当前可执行文件所在路径
	s := file.GetCurrentPath()

	// 设置 JDK 存储目录路径
	cfx.Store = filepath.Join(s, "store")

	// 设置下载临时目录路径
	cfx.Download = filepath.Join(s, "download")

	// 如果未配置原始路径，使用默认值
	if cfx.Originalpath == "" {
		cfx.Originalpath = defaultOriginalpath
	}

	// 如果配置了代理，设置 HTTP 代理
	if cfx.Proxy != "" {
		web.SetProxy(cfx.Proxy)
	}

	return nil
}

// shutdown 在应用关闭后执行
// 主要功能：保存配置到 jvms.json 文件
// 参数:
//   c - CLI 上下文
// 返回:
//   error - 保存配置失败时返回错误
func shutdown(c *cli.Context) error {
	if err := store.Save("jvms.json", &cfx); err != nil {
		return errors.New("failed to save the config:" + err.Error())
	}
	return nil
}
