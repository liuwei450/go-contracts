package util

import (
	"github.com/ethereum/go-ethereum/log"
	"os"
)

// 全局日志实例
var Log = log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, log.LevelInfo, true))

// InitLogger 初始化日志（可通过外部配置调整级别）
func InitLogger() {
	// 后续可扩展：从配置文件/命令行参数读取日志级别
	// 示例：logLevel := log.LevelInfo; if debug { logLevel = log.LevelDebug }
	// Log = log.NewLogger(log.NewTerminalHandlerWithLevel(os.Stderr, logLevel, true))
}
