package synchronizer

import (
	"context"
	"errors"
	"fmt"
	"go-contracts/config"
	"go-contracts/models"
	"go-contracts/synchronizer/node"
	"go-contracts/util"
	"os"
	"sync/atomic"
	"time"
)

// BlockBatch 表示一批区块数据
type BlockBatch struct {
	Blocks []*models.Block
}

// Synchronizer 从区块链同步事件的组件
type Synchronizer struct {
	interval     time.Duration           // 同步间隔
	shutdown     context.CancelCauseFunc // 取消函数（用于主动退出）
	stopped      atomic.Bool             // 停止状态标记
	ethClient    node.EthClient          // 以太坊客户端
	blockChannel chan<- *BlockBatch      // 区块数据通道
	lastBlockNum uint64                  // 最后处理的区块号
}

// NewSynchronizer 创建同步器实例
func NewSynchronizer(cfg *config.IndexerConfig, ethClient node.EthClient, blockChannel chan<- *BlockBatch, shutdown context.CancelCauseFunc) (*Synchronizer, error) {
	// 从配置读取同步间隔（默认 10 秒）
	interval := time.Duration(cfg.Interval) * time.Second
	if interval <= 0 {
		interval = 10 * time.Second
	}

	// 初始化最后处理的区块号
	// 在实际环境中，应该从数据库读取最后处理的区块号
	// 这里设置为100，跳过已经存在的区块记录
	lastBlockNum := uint64(100)
	
	return &Synchronizer{
		interval:     interval,
		shutdown:     shutdown,
		ethClient:    ethClient,
		blockChannel: blockChannel,
		lastBlockNum: lastBlockNum,
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
	
	// 模拟获取最新区块号
	// 实际实现应该通过 ethClient 获取最新区块号
	// 为了演示，这里简单模拟
	startBlock := s.lastBlockNum
	endBlock := startBlock + 19 // 每次扫描20个区块（包括startBlock）
	
	util.Log.Info("开始扫描区块范围", "start", startBlock, "end", endBlock)
	
	// 创建区块批次
	blockBatch := &BlockBatch{
		Blocks: make([]*models.Block, 0, 20),
	}
	
	// 模拟扫描20个区块
	// 实际实现应该使用 ethClient 获取每个区块的数据
	for i := startBlock; i <= endBlock; i++ {
		// 模拟区块数据
		// 在实际实现中，这里应该调用 ethclient 获取区块数据，然后转换为 models.Block
		// 生成固定长度的模拟区块哈希，确保不会超过数据库字段长度限制
		// 使用i的哈希值加上随机字符串，然后截取固定长度
		blockNumHex := fmt.Sprintf("%x", i)
		parentNumHex := fmt.Sprintf("%x", i-1)
		blockHash := fmt.Sprintf("0x%064s", blockNumHex)[:66]  // 确保长度为66
		parentHash := fmt.Sprintf("0x%064s", parentNumHex)[:66]  // 确保长度为66
		
		// 从环境变量读取矿工地址，使用默认值作为备选
		minerAddress := os.Getenv("MINER_ADDRESS")
		if minerAddress == "" {
			return errors.New("MINER_ADDRESS 环境变量未设置")
		}
		
		block := &models.Block{
			BlockNumber:  i,
			BlockHash:    blockHash,
			ParentHash:   parentHash,
			TxRoot:       "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",// 交易根哈希
			ReceiptsRoot: "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",// 收据根哈希
			StateRoot:    "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",// 状态根哈希
			Miner:        minerAddress,// 矿工地址（从环境变量读取）
			GasUsed:      21000 * uint64(i%10+1),
			GasLimit:     10000000,
			Time:         uint64(time.Now().Unix() - int64(i)),
			Timestamp:    time.Unix(int64(uint64(time.Now().Unix())-i), 0),
			ExtraData:    "0x",
			Transactions: uint(i % 20),
		}
		blockBatch.Blocks = append(blockBatch.Blocks, block)
	}
	
	// 发送区块批次到通道
	select {
	case s.blockChannel <- blockBatch:
		util.Log.Info("成功发送区块批次到处理通道", "count", len(blockBatch.Blocks))
		s.lastBlockNum = endBlock + 1 // 更新最后处理的区块号
	case <-ctx.Done():
		return ctx.Err()
	}
	
	return nil
}

// Close 停止同步器
func (s *Synchronizer) Close() error {
	if s.stopped.CompareAndSwap(false, true) {
		util.Log.Info("同步器已停止")
	}
	return nil
}
