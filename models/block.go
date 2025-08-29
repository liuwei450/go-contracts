package models

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Block 区块实体定义区块链区块的基本信息
// 该结构用于存储和处理从区块链同步的区块数据
// 对应数据库中的blocks表（如果存在）
type Block struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BlockNumber   uint64    `gorm:"uniqueIndex" json:"block_number"`         // 区块号，唯一索引
	BlockHash     string    `gorm:"size:66;index" json:"block_hash"`        // 区块哈希，索引
	ParentHash    string    `gorm:"size:66" json:"parent_hash"`             // 父区块哈希
	TxRoot        string    `gorm:"size:66" json:"tx_root"`                 // 交易树根哈希
	ReceiptsRoot  string    `gorm:"size:66" json:"receipts_root"`           // 收据树根哈希
	StateRoot     string    `gorm:"size:66" json:"state_root"`              // 状态树根哈希
	Miner         string    `gorm:"size:42;index" json:"miner"`             // 矿工地址，索引
	GasUsed       uint64    `json:"gas_used"`                                // 已用Gas
	GasLimit      uint64    `json:"gas_limit"`                               // Gas上限
	Time          uint64    `json:"time"`                                    // 区块时间戳（Unix时间）
	Timestamp     time.Time `gorm:"index" json:"timestamp"`                 // 格式化的时间戳，索引
	ExtraData     string    `json:"extra_data"`                               // 额外数据
	Transactions  uint      `json:"transactions"`                             // 交易数量
	CreatedAt     time.Time `json:"created_at"`                               // 记录创建时间
	UpdatedAt     time.Time `json:"updated_at"`                               // 记录更新时间
}

// NewBlockFromRPC 从RPC响应数据创建区块实体
// 这个方法用于将从以太坊节点获取的区块数据转换为我们的模型结构
func NewBlockFromRPC(blockNumber uint64, blockHash common.Hash, parentHash common.Hash, 
	txRoot common.Hash, receiptsRoot common.Hash, stateRoot common.Hash, 
	miner common.Address, gasUsed, gasLimit uint64, timestamp uint64, 
	extraData []byte, txCount int) *Block {
	return &Block{
		BlockNumber:  blockNumber,
		BlockHash:    blockHash.Hex(),
		ParentHash:   parentHash.Hex(),
		TxRoot:       txRoot.Hex(),
		ReceiptsRoot: receiptsRoot.Hex(),
		StateRoot:    stateRoot.Hex(),
		Miner:        miner.Hex(),
		GasUsed:      gasUsed,
		GasLimit:     gasLimit,
		Time:         timestamp,
		Timestamp:    time.Unix(int64(timestamp), 0),
		ExtraData:    hexutil.Encode(extraData),
		Transactions: uint(txCount),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}