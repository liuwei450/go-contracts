package models

import (
	"time"
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// AirdropEvent 存储空投事件的数据模型
// 同时支持ERC20和BNB的空投事件

type AirdropEvent struct {
	gorm.Model
	
	// 基础交易信息
	TransactionHash common.Hash `gorm:"size:66;index" json:"transaction_hash"` // 交易哈希
	BlockNumber     uint64       `gorm:"index" json:"block_number"`             // 区块号
	BlockTime       time.Time    `json:"block_time"`                             // 区块时间
	
	// 事件特有信息
	EventType       string       `gorm:"size:50;index" json:"event_type"`       // 事件类型：AirdropERC20 或 AirdropBNB
	Recipient       common.Address `gorm:"size:42;index" json:"recipient"`      // 接收者地址
	Amount          string       `gorm:"size:100" json:"amount"`                // 金额（以字符串形式存储大整数）
	TokenAddress    common.Address `gorm:"size:42;index" json:"token_address"`  // 代币地址（BNB空投时为0x0000000000000000000000000000000000000000）
	ContractAddress common.Address `gorm:"size:42;index" json:"contract_address"` // 空投合约地址
}

// TableName 自定义表名
func (AirdropEvent) TableName() string {
	return "airdrop_events"
}