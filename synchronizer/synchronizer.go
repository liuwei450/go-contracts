package synchronizer

import (
	"context"
	"go-contracts/config"
	"go-contracts/util"
	"sync/atomic"
	"time"
)

// Synchronizer 从区块链同步事件的组件
type Synchronizer struct {
	interval time.Duration           // 同步间隔
	shutdown context.CancelCauseFunc // 取消函数（用于主动退出）
	stopped  atomic.Bool             // 停止状态标记
}

// NewSynchronizer 创建同步器实例
func NewSynchronizer(cfg *config.IndexerConfig, shutdown context.CancelCauseFunc) (*Synchronizer, error) {
	// 从配置读取同步间隔（默认 10 秒）
	interval := time.Duration(cfg.Interval) * time.Second
	if interval <= 0 {
		interval = 10 * time.Second
	}

	return &Synchronizer{
		interval: interval,
		shutdown: shutdown,
	}, nil
}

// Start 启动同步器（非阻塞）
func (s *Synchronizer) Start(ctx context.Context) error {
	if s.stopped.Load() {
		return nil
	}

	util.Log.Info("同步器启动", "interval", s.interval)
	go s.runLoop(ctx) // 启动同步循环（后台协程）
	return nil
}

// runLoop 定时执行同步任务
func (s *Synchronizer) runLoop(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.stopped.Store(true)
			util.Log.Info("同步器退出循环")
			return
		case <-ticker.C:
			if err := s.syncOnce(ctx); err != nil {
				util.Log.Error("同步任务失败", "err", err)
				s.shutdown(err) // 同步失败时主动触发服务退出
			}
		}
	}
}

// syncOnce 单次同步逻辑（实际从区块链拉取事件）
func (s *Synchronizer) syncOnce(ctx context.Context) error {
	util.Log.Debug("执行区块链事件同步...")
	// 示例：调用区块链 RPC 获取最新区块/事件
	// 实际实现：使用 ethclient 获取日志，处理后交给 Processor
	return nil
}

// Close 停止同步器
func (s *Synchronizer) Close() error {
	if s.stopped.CompareAndSwap(false, true) {
		util.Log.Info("同步器已停止")
	}
	return nil
}
