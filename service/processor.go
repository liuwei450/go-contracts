package service

import (
	"context"
	"go-contracts/config"
	"go-contracts/database"
	"go-contracts/synchronizer"
	"go-contracts/util"
	"sync/atomic"
)

// Processor 处理同步后的数据（如入库、过滤、转换）
type Processor struct {
	db          *database.DB            // 数据库连接
	shutdown    context.CancelCauseFunc // 取消函数
	stopped     atomic.Bool             // 停止状态标记
	blockChannel <-chan *synchronizer.BlockBatch // 区块数据通道
}

// NewProcessor 创建处理器实例
func NewProcessor(cfg *config.IndexerConfig, db *database.DB, blockChannel <-chan *synchronizer.BlockBatch, shutdown context.CancelCauseFunc) (*Processor, error) {
	return &Processor{
		db:          db,
		shutdown:    shutdown,
		blockChannel: blockChannel,
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

// processLoop 处理同步后的数据（从 channel 接收区块数据并保存到数据库）
func (p *Processor) processLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			p.stopped.Store(true)
			util.Log.Info("处理器退出循环")
			return
		case blockBatch, ok := <-p.blockChannel:
			if !ok {
				util.Log.Info("区块通道已关闭，处理器退出循环")
				p.stopped.Store(true)
				return
			}
			
			util.Log.Info("接收到区块批次", "count", len(blockBatch.Blocks))
			
			// 批量保存区块数据到数据库
			if err := p.db.Create(blockBatch.Blocks).Error; err != nil {
				util.Log.Error("区块数据保存失败", "err", err)
				p.shutdown(err)
				return
			}
			
			util.Log.Info("区块数据保存成功", "count", len(blockBatch.Blocks))
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
