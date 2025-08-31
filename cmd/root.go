package cmd

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
	"go-contracts/config"
	"go-contracts/controller"
	"go-contracts/cycle"
	"go-contracts/database"
	"go-contracts/service"
	"go-contracts/util"
	"os"
	"os/signal"
	"syscall"
)

// 全局命令行参数（示例）
var globalFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "配置文件路径",
		Value:   "config.yaml",
	},
	&cli.BoolFlag{
		Name:    "debug",
		Aliases: []string{"d"},
		Usage:   "开启调试日志",
	},
}

func runIndexer(ctx *cli.Context, shutdown context.CancelCauseFunc) (cycle.Service, error) {
	log.Info("run event watcher indexer")
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		log.Error("failed to load config", "err", err)
		return nil, err
	}
	return service.NewIndexerService(ctx, cfg, shutdown)
}

func runAirdropWatcher(ctx *cli.Context, shutdown context.CancelCauseFunc) (cycle.Service, error) {
	log.Info("run airdrop event watcher")
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		log.Error("failed to load config", "err", err)
		return nil, err
	}
	return service.NewAirdropWatcher(ctx, cfg, shutdown)
}

func runApi(ctx *cli.Context, _ context.CancelCauseFunc) (cycle.Service, error) {
	log.Info("running api...")
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		log.Error("failed to load config", "err", err)
		return nil, err
	}

	return controller.NewApi(ctx, cfg)
}

func StartServer(gitCommit, gitDate string) *cli.App {
	return &cli.App{
		Name:                 "event-indexer",
		Version:              versionWithCommit(gitCommit, gitDate),
		Description:          "Optimism 事件索引器（带 API 服务）",
		EnableBashCompletion: true,
		Flags:                globalFlags,
		Commands: []*cli.Command{
			{
				Name:        "api",
				Usage:       "启动 API 服务",
				Description: "提供 HTTP 接口查询索引数据",
				Flags: append(globalFlags, []cli.Flag{ // 服务专属参数
					&cli.StringFlag{
						Name:  "port",
						Usage: "API 服务端口",
						Value: "8080",
					},
				}...),
				Action: cycle.LifecycleCmd(runApi), // 绑定 API 服务
			},
			{
				Name:        "index",
				Usage:       "启动索引服务",
				Description: "从区块链同步事件并索引到数据库",
				Flags: append(globalFlags, []cli.Flag{ // 服务专属参数
					&cli.IntFlag{
						Name:  "interval",
						Usage: "索引任务间隔（秒）",
						Value: 10,
					},
				}...),
				Action: cycle.LifecycleCmd(runIndexer), // 绑定索引服务
		},
		{
			Name:        "airdrop-watch",
			Usage:       "启动空投事件监听服务",
			Description: "监听空投合约的AirdropERC20和AirdropBNB事件并保存到数据库",
			Flags:       globalFlags,
			Action:      cycle.LifecycleCmd(runAirdropWatcher), // 绑定空投监听服务
		},
		{
			Name:        "migrate",
				Usage:       "执行数据库迁移",
				Description: "初始化或更新数据库表结构",
				Flags:       globalFlags,
				Action:      runMigrations, // 一次性任务（无需优雅关停）
			},
		},
	}
}

// versionWithCommit 生成版本信息（包含 Git 提交和日期）
func versionWithCommit(commit, date string) string {
	version := "v0.1.0"
	if commit != "" {
		version += fmt.Sprintf(" (commit: %s, date: %s)", commit, date)
	}
	return version
}

// runMigrations 数据库迁移
func runMigrations(ctx *cli.Context) error {
	util.Log.Info("执行数据库迁移...")

	// 1. 加载配置
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		util.Log.Error("加载配置失败", "err", err)
		return fmt.Errorf("load config: %w", err) // 使用 %w 包装错误，保留调用栈
	}

	// 2. 创建带中断监听的上下文
	c, cancel := context.WithCancel(context.Background())
	defer cancel() // 确保函数退出时释放上下文

	// 监听系统中断信号（如 Ctrl+C），触发迁移中断
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		util.Log.Warn("收到中断信号，正在取消迁移...")
		cancel() // 取消上下文，终止后续操作
	}()

	// 3. 执行 SQL 迁移（基于 golang-migrate 或其他工具）
	if err := database.ExecuteMigrations(c, cfg.MasterDB.DSN(), cfg.MasterDB.Driver, cfg.MigrationDir); err != nil {
		util.Log.Error("迁移执行失败", "err", err, "migration_dir", cfg.MigrationDir)
		return fmt.Errorf("execute migration: %w", err)
	}
	// 示例：调用 gorm-migrate 或其他迁移工具
	util.Log.Info("数据库迁移完成")
	return nil
}
