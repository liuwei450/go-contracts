package service

import (
	"context"
	"math/big"
	"time"
	"sync/atomic"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
	"go-contracts/config"
	"go-contracts/contract"
	"go-contracts/database"
	"go-contracts/models"
	"go-contracts/util"
)

// AirdropWatcher 空投事件监听服务
// 负责监听和处理AirdropERC20和AirdropBNB事件
// 实现了cycle.Service接口

type AirdropWatcher struct {
	shutdown    context.CancelCauseFunc  // 取消函数
	stopped     atomic.Bool              // 停止状态标记
	db          *database.DB             // 数据库连接
	ethClient   *ethclient.Client        // 以太坊客户端
	contract    *contract.Airdrop        // 空投合约实例
	contractAddr common.Address          // 空投合约地址
}

// Stopped 实现cycle.Service接口，返回服务是否已停止
func (w *AirdropWatcher) Stopped() bool {
	return w.stopped.Load()
}

// NewAirdropWatcher 创建空投事件监听服务实例
func NewAirdropWatcher(c *cli.Context, cfg *config.Config, shutdown context.CancelCauseFunc) (*AirdropWatcher, error) {
	// 初始化数据库连接
	db, err := database.NewDb(c.Context, &cfg.MasterDB)
	if err != nil {
		util.Log.Error("初始化数据库失败", "err", err)
		return nil, err
	}

	// 初始化以太坊客户端
	rawURL := config.RAW_URL
	ethClient, err := ethclient.Dial(rawURL)
	if err != nil {
		util.Log.Error("连接BSC节点失败", "url", rawURL, "err", err)
		db.Close()
		return nil, err
	}

	// 初始化空投合约实例
	contractAddr := common.HexToAddress(config.AIRDROP_CONTRACT_ADDRESS)
	airdropContract, err := contract.NewAirdrop(contractAddr, ethClient)
	if err != nil {
		util.Log.Error("初始化空投合约失败", "addr", contractAddr, "err", err)
		db.Close()
		ethClient.Close()
		return nil, err
	}

	return &AirdropWatcher{
		shutdown:    shutdown,
		db:          db,
		ethClient:   ethClient,
		contract:    airdropContract,
		contractAddr: contractAddr,
	},
	nil
}

// Start 启动监听服务
func (w *AirdropWatcher) Start(ctx context.Context) error {
	if w.stopped.Load() {
		return nil
	}

	util.Log.Info("空投事件监听服务启动", "contract", w.contractAddr.Hex())
	
	// 启动两个事件监听协程
	go w.watchAirdropERC20(ctx)
	go w.watchAirdropBNB(ctx)
	
	return nil
}

// Stop 停止监听服务
func (w *AirdropWatcher) Stop(ctx context.Context) error {
	if w.stopped.CompareAndSwap(false, true) {
		util.Log.Info("空投事件监听服务停止中...")
		
		// 关闭资源
		if w.db != nil {
			w.db.Close()
		}
		
		if w.ethClient != nil {
			w.ethClient.Close()
		}
		
		util.Log.Info("空投事件监听服务已停止")
	}
	return nil
}

// watchAirdropERC20 监听AirdropERC20事件
func (w *AirdropWatcher) watchAirdropERC20(ctx context.Context) {
	// 创建事件过滤器
	query := &bind.WatchOpts{
		Context: ctx,
	}
	
	// 创建事件接收通道
	logs := make(chan *contract.AirdropAirdropERC20)
	
	// 监听事件
	sub, err := w.contract.WatchAirdropERC20(query, logs, []common.Address{})
	if err != nil {
		util.Log.Error("监听AirdropERC20事件失败", "err", err)
		w.shutdown(err)
		return
	}
	defer sub.Unsubscribe()
	
	util.Log.Info("开始监听AirdropERC20事件")
	
	// 处理事件流
	for {
		select {
		case <-ctx.Done():
			util.Log.Info("AirdropERC20事件监听停止")
			return
		case err := <-sub.Err():
			util.Log.Error("AirdropERC20事件订阅错误", "err", err)
			// 不立即终止服务，尝试重新连接
			util.Log.Info("尝试重新监听AirdropERC20事件")
			go w.watchAirdropERC20(ctx)
			return
		case event := <-logs:
			// 处理单个事件
			w.handleAirdropEvent(ctx, "AirdropERC20", event)
		}
	}
}

// watchAirdropBNB 监听AirdropBNB事件
func (w *AirdropWatcher) watchAirdropBNB(ctx context.Context) {
	// 创建事件过滤器
	query := &bind.WatchOpts{
		Context: ctx,
	}
	
	// 创建事件接收通道
	logs := make(chan *contract.AirdropAirdropBNB)
	
	// 监听事件
	sub, err := w.contract.WatchAirdropBNB(query, logs, []common.Address{})
	if err != nil {
		util.Log.Error("监听AirdropBNB事件失败", "err", err)
		w.shutdown(err)
		return
	}
	defer sub.Unsubscribe()
	
	util.Log.Info("开始监听AirdropBNB事件")
	
	// 处理事件流
	for {
		select {
		case <-ctx.Done():
			util.Log.Info("AirdropBNB事件监听停止")
			return
		case err := <-sub.Err():
			util.Log.Error("AirdropBNB事件订阅错误", "err", err)
			// 不立即终止服务，尝试重新连接
			util.Log.Info("尝试重新监听AirdropBNB事件")
			go w.watchAirdropBNB(ctx)
			return
		case event := <-logs:
			// 处理单个事件
			w.handleAirdropEvent(ctx, "AirdropBNB", event)
		}
	}
}

// handleAirdropEvent 处理空投事件
func (w *AirdropWatcher) handleAirdropEvent(ctx context.Context, eventType string, event interface{}) {
	var recipient common.Address
	var amount *big.Int
	var rawLog types.Log
	
	switch e := event.(type) {
	case *contract.AirdropAirdropERC20:
		recipient = e.Recipient
		amount = e.Amount
		rawLog = e.Raw
	case *contract.AirdropAirdropBNB:
		recipient = e.Recipient
		amount = e.Amount
		rawLog = e.Raw
	default:
		util.Log.Error("未知的事件类型", "type", eventType)
		return
	}
	
	// 获取区块信息
	block, err := w.ethClient.BlockByHash(ctx, rawLog.BlockHash)
	if err != nil {
		util.Log.Error("获取区块信息失败", "hash", rawLog.BlockHash.Hex(), "err", err)
		return
	}
	
	// 创建事件记录
	dbEvent := &models.AirdropEvent{
		TransactionHash: rawLog.TxHash,
		BlockNumber:     block.NumberU64(),
		BlockTime:       time.Unix(int64(block.Time()), 0),
		EventType:       eventType,
		Recipient:       recipient,
		Amount:          amount.String(),
		ContractAddress: w.contractAddr,
	}
	
	// 根据事件类型设置TokenAddress
	if eventType == "AirdropBNB" {
		dbEvent.TokenAddress = common.HexToAddress(config.NATIVE_TOKEN_ADDRESS)
	} else {
		// 获取代币地址
		tokenAddr, err := w.contract.Token(&bind.CallOpts{Context: ctx})
		if err == nil {
			dbEvent.TokenAddress = tokenAddr
		} else {
			util.Log.Warn("获取代币地址失败", "err", err)
		}
	}
	
	// 保存到数据库
	if err := w.db.Create(dbEvent).Error;
	err != nil {
		util.Log.Error("保存空投事件失败", "err", err, "event", dbEvent)
	} else {
		util.Log.Info("空投事件保存成功", "type", eventType, "recipient", recipient.Hex(), "amount", amount.String())
	}
}