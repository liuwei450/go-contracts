package service

import (
	"context"
	"go-contracts/config"
	"go-contracts/database"
	"go-contracts/util"
	"sync/atomic"
)

// Processor 处理同步后的数据（如入库、过滤、转换）
type Processor struct {
	db       *database.DB            // 数据库连接
	shutdown context.CancelCauseFunc // 取消函数
	stopped  atomic.Bool             // 停止状态标记
}

// NewProcessor 创建处理器实例
func NewProcessor(cfg *config.IndexerConfig, db *database.DB, shutdown context.CancelCauseFunc) (*Processor, error) {
	return &Processor{
		db:       db,
		shutdown: shutdown,
	}, nil
}

// Start 启动处理器（非阻塞）
func (p *Processor) Start(ctx context.Context) error {
	if p.stopped.Load() {
		return nil
	}

	util.Log.Info("数据处理器启动")
	// 示例：启动消息队列消费者或处理协程
	go p.processLoop(ctx)
	return nil
}

// processLoop 处理同步后的数据（示例逻辑）
func (p *Processor) processLoop(ctx context.Context) {
	// 实际实现：从 channel 接收 Synchronizer 推送的事件，处理后入库
	for {
		select {
		case <-ctx.Done():
			p.stopped.Store(true)
			util.Log.Info("处理器退出循环")
			return
			// case event := <-eventCh: // 假设从同步器接收事件
			// 	if err := p.db.Save(event).Error; err != nil {
			// 		util.Log.Error("事件入库失败", "err", err)
			// 		p.shutdown(err)
			// 	}
		}
	}
}

// Close 停止处理器
func (p *Processor) Close() error {
	if p.stopped.CompareAndSwap(false, true) {
		util.Log.Info("数据处理器已停止")
	}
	return nil
}
