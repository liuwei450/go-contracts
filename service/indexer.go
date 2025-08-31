package service

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
	"go-contracts/config"
	"go-contracts/contract"
	"go-contracts/database"
	"go-contracts/synchronizer"
	"go-contracts/synchronizer/node"
	"go-contracts/util"
	"math/big"
	"sync/atomic"
	"time"
)

// MockEthClient 模拟的以太坊客户端实现，用于开发和测试
// 实际项目中应该使用真实的以太坊节点连接

type MockEthClient struct {}

// 实现 node.EthClient 接口的方法（简化版本）
func (m *MockEthClient) GetERC20Contract(ctx context.Context, contractAddress common.Address) (*contract.Erc20, error) {
	// 实际实现应该返回真实的合约实例
	return nil, nil
}

func (m *MockEthClient) ERC20Allowance(ctx context.Context, contractAddress common.Address, owner, spender common.Address) (*big.Int, error) {
	// 实际实现应该查询真实的授权额度
	return big.NewInt(0), nil
}

func (m *MockEthClient) ERC20Approve(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, spender common.Address, value *big.Int) (*big.Int, error) {
	// 实际实现应该执行真实的授权操作
	return big.NewInt(0), nil
}

func (m *MockEthClient) ERC20Transfer(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, to common.Address, value *big.Int) (*big.Int, error) {
	// 实际实现应该执行真实的转账操作
	return big.NewInt(0), nil
}

func (m *MockEthClient) ERC20TransferFrom(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, from, to common.Address, value *big.Int) (*big.Int, error) {
	// 实际实现应该执行真实的授权转账操作
	return big.NewInt(0), nil
}

func (m *MockEthClient) ERC20Balance(ctx context.Context, contractAddress common.Address, account common.Address) (*big.Int, error) {
	// 实际实现应该查询真实的余额
	return big.NewInt(0), nil
}

func (m *MockEthClient) ERC20TotalSupply(ctx context.Context, contractAddress common.Address) (*big.Int, error) {
	// 实际实现应该查询真实的总供应量
	return big.NewInt(0), nil
}

func (m *MockEthClient) ERC20TokenInfo(ctx context.Context, contractAddress common.Address) (string, string, uint8, error) {
	// 实际实现应该查询真实的代币信息
	return "MockToken", "MOCK", 18, nil
}

// 事件索引服务实现 cli.Service 接口
type IndexerService struct {
	ticker      *time.Ticker             // 定时索引任务
	shutdown    context.CancelCauseFunc  // 取消函数
	stopped     atomic.Bool              // 停止状态标记
	db          *database.DB             // 数据库连接
	synchronizer *synchronizer.Synchronizer // 同步器
	processor   *Processor               // 处理器
	blockChannel chan *synchronizer.BlockBatch // 区块数据通道
}

// 创建索引服务实例（业务逻辑入口）
func NewIndexerService(c *cli.Context, cfg *config.Config, shutdown context.CancelCauseFunc) (*IndexerService, error) {
	// 1. 从配置读取索引间隔
	interval := cfg.Indexer.Interval
	if interval <= 0 {
		interval = 10 // 默认 10 秒
	}
	
	// 2. 创建区块数据通道
	blockChannel := make(chan *synchronizer.BlockBatch, 100)
	
	// 3. 初始化外部依赖：数据库连接
	db, err := database.NewDb(c.Context, &cfg.MasterDB)
	if err != nil {
		util.Log.Error("初始化数据库失败", "err", err)
		return nil, err
	}
	
	// 4. 初始化外部依赖：区块链客户端
	// 注意：实际项目中应该从配置中读取RPC URL
	// 这里使用模拟客户端，因为我们没有实际的RPC URL配置
	// 实际实现应该是：ethClient, err := node.DialEthClient(*c, cfg.EthRPCUrl)
	var ethClient node.EthClient = &MockEthClient{}
	
	// 5. 创建核心组件：同步器（从区块链拉取事件）
	sync, err := synchronizer.NewSynchronizer(&cfg.Indexer, ethClient, blockChannel, shutdown)
	if err != nil {
		util.Log.Error("初始化同步器失败", "err", err)
		db.Close()
		close(blockChannel)
		return nil, err
	}
	
	// 6. 创建核心组件：处理器（处理事件并入库）
	processor, err := NewProcessor(&cfg.Indexer, db, blockChannel, shutdown)
	if err != nil {
		util.Log.Error("初始化处理器失败", "err", err)
		db.Close()
		close(blockChannel)
		return nil, err
	}
	
	// 7. 组装索引服务实例（包含所有组件和依赖）
	service := &IndexerService{
		ticker:      time.NewTicker(time.Duration(interval) * time.Second),
		shutdown:    shutdown,
		db:          db,
		synchronizer: sync,
		processor:   processor,
		blockChannel: blockChannel,
	}
	
	return service, nil
}

// Start 启动索引服务（阻塞方法）
func (s *IndexerService) Start(ctx context.Context) error {
	util.Log.Info("索引服务启动，开始监听事件")

	// 启动处理器
	if err := s.processor.Start(ctx); err != nil {
		util.Log.Error("启动处理器失败", "err", err)
		return err
	}

	// 启动同步器
	if err := s.synchronizer.Start(ctx); err != nil {
		util.Log.Error("启动同步器失败", "err", err)
		// 确保处理器也停止
		s.processor.Close()
		return err
	}

	// 保持服务运行，直到收到停止信号
	<-ctx.Done()
	return ctx.Err()
}

// 停止服务（清理资源）
func (s *IndexerService) Stop(ctx context.Context) error {
	util.Log.Info("索引服务清理资源...")
	
	// 标记服务为已停止
	s.stopped.Store(true)
	
	// 停止同步器
	if s.synchronizer != nil {
		s.synchronizer.Close()
	}
	
	// 停止处理器
	if s.processor != nil {
		s.processor.Close()
	}
	
	// 关闭区块通道
	if s.blockChannel != nil {
		close(s.blockChannel)
	}
	
	// 关闭数据库连接
	if s.db != nil {
		s.db.Close()
	}
	
	// 停止定时器（如果仍在使用）
	if s.ticker != nil {
		s.ticker.Stop()
	}
	
	util.Log.Info("索引服务已成功停止并清理所有资源")
	return nil
}

// Stopped 返回服务是否已停止
func (s *IndexerService) Stopped() bool {
	return s.stopped.Load()
}
