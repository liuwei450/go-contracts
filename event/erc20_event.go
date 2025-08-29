package event

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

// ERC20ApprovalEvent 表示ERC20授权事件
type ERC20ApprovalEvent struct {
	Owner    common.Address // 授权方地址
	Spender  common.Address // 被授权方地址
	Value    *big.Int       // 授权金额
	TxHash   common.Hash    // 交易哈希
	BlockNumber *big.Int    // 区块号
}

// ERC20TransferEvent 表示ERC20转账事件
type ERC20TransferEvent struct {
	From     common.Address // 发送方地址
	To       common.Address // 接收方地址
	Value    *big.Int       // 转账金额
	TxHash   common.Hash    // 交易哈希
	BlockNumber *big.Int    // 区块号
}

// ERC20LogData 表示ERC20事件日志数据
type ERC20LogData struct {
	Log      types.Log      // 原始日志数据
	EventName string        // 事件名称
	BlockNumber *big.Int    // 区块号
	TxHash   common.Hash    // 交易哈希
	ContractAddress common.Address // 合约地址
}