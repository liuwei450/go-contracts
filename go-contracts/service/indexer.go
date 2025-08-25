package service

import (
	"context"
	"github.com/urfave/cli/v2"
	"go-contracts/cycle"
	"go-contracts/util"
	"time"
)

// 事件索引服务实现 cli.Service 接口
type IndexerService struct {
	ticker *time.Ticker // 定时索引任务
}

// 创建索引服务实例（业务逻辑入口）
func NewIndexerService(c *cli.Context) (cycle.Service, error) {
	// 示例：从命令行参数读取索引间隔
	interval := c.Int("interval")
	if interval <= 0 {
		interval = 10 // 默认 10 秒
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
func (s *IndexerService) Stop() error {
	util.Log.Info("索引服务清理资源...")
	s.ticker.Stop() // 停止定时器
	// 示例：关闭数据库连接、保存进度等
	return nil
}
