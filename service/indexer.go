package service

import (
	"context"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
	"go-contracts/database"

	"go-contracts/config"
	"sync/atomic"

	"go-contracts/util"
	"time"
)

// 事件索引服务实现 cli.Service 接口
type IndexerService struct {
	ticker   *time.Ticker // 定时索引任务
	shutdown context.CancelCauseFunc
	stopped  atomic.Bool
}

// 创建索引服务实例（业务逻辑入口）
func NewIndexerService(c *cli.Context, cfg *config.Config, shutdown context.CancelCauseFunc) (*IndexerService, error) {
	// 示例：从命令行参数读取索引间隔
	interval := c.Int("interval")
	if interval <= 0 {
		interval = 10 // 默认 10 秒
	}
	//1.初始化外部依赖：区块链客户端

	//2.初始化外部依赖：数据库连接
	_, err := database.NewDb(c.Context, &cfg.MasterDB)
	if err != nil {
		log.Error("初始化数据库失败", err)
		return nil, err
	}
	// 3. 创建核心组件：同步器（从区块链拉取事件）

	//4.创建核心组件：处理器（处理事件并入库）
	//5. 组装索引服务实例（包含所有组件和依赖）
	if err != nil {
		return nil, err
	}
	return &IndexerService{
		ticker: time.NewTicker(time.Duration(interval) * time.Second),
	}, nil
}

// Start 启动索引服务（阻塞方法）
func (s *IndexerService) Start(ctx context.Context) error {
	util.Log.Info("索引服务启动，开始监听事件")

	// 定时执行索引任务
	for {
		select {
		case <-ctx.Done():
			return ctx.Err() // 收到退出信号，停止服务
		case <-s.ticker.C:
			s.runIndexTask() // 执行索引任务
		}
	}
}

// 模拟索引任务（实际业务逻辑）
func (s *IndexerService) runIndexTask() {
	util.Log.Debug("执行事件索引...")
	// 示例：查询区块链事件、写入数据库等
}

// 停止服务（清理资源）
func (s *IndexerService) Stop(ctx context.Context) error {
	util.Log.Info("索引服务清理资源...")
	s.ticker.Stop() // 停止定时器
	// 示例：关闭数据库连接、保存进度等
	return nil
}

// Stopped 返回服务是否已停止
func (s *IndexerService) Stopped() bool {
	return s.stopped.Load()
}
