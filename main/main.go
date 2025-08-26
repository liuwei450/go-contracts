package main

import (
	"go-contracts/cmd"
	"go-contracts/util"
	"os"
)

var (
	GitCommit = "dev" // 例如：git rev-parse --short HEAD
	GitData   = "now" // 例如：date -u +%Y-%m-%dT%H:%M:%SZ
)

func main() {
	// 初始化全局日志
	util.InitLogger()

	// 创建 CLI 应用并运行
	app := cmd.StartServer(GitCommit, GitData)
	if err := app.Run(os.Args); err != nil {
		util.Log.Crit("应用启动失败", "error", err)
		os.Exit(1)
	}

}
