package models

import (
	"time"
)

// ERC20TokenInfo 表示ERC20代币的基本信息
type ERC20TokenInfo struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Address     string         `gorm:"size:42;uniqueIndex" json:"address"` // 合约地址
	Name        string         `gorm:"size:100" json:"name"`               // 代币名称
	Symbol      string         `gorm:"size:50" json:"symbol"`              // 代币符号
	Decimals    uint8          `json:"decimals"`                            // 小数位数
	TotalSupply string         `gorm:"type:text" json:"total_supply"`     // 总供应量
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ERC20Transaction 表示ERC20代币交易记录
type ERC20Transaction struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	TxHash          string         `gorm:"size:66;index" json:"tx_hash"`      // 交易哈希
	BlockHash       string         `gorm:"size:66;index" json:"block_hash"`   // 区块哈希
	BlockNumber     uint64         `gorm:"index" json:"block_number"`          // 区块号
	From            string         `gorm:"size:42;index" json:"from"`          // 发送方地址
	To              string         `gorm:"size:42;index" json:"to"`            // 接收方地址
	ContractAddress string         `gorm:"size:42;index" json:"contract_address"` // 合约地址
	Amount          string         `gorm:"type:text" json:"amount"`           // 交易金额
	GasUsed         uint64         `json:"gas_used"`                            // 消耗的Gas
	GasPrice        string         `gorm:"type:text" json:"gas_price"`        // Gas价格
	Status          bool           `json:"status"`                              // 交易状态
	TransactionType string         `gorm:"size:50" json:"transaction_type"`   // 交易类型：transfer, transfer_from, approve
	CreatedAt       time.Time      `json:"created_at"`
}

// ERC20Balance 表示用户的ERC20代币余额
type ERC20Balance struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ContractAddress string         `gorm:"size:42;index:idx_contract_address_account" json:"contract_address"` // 合约地址
	Account         string         `gorm:"size:42;index:idx_contract_address_account" json:"account"`         // 用户地址
	Balance         string         `gorm:"type:text" json:"balance"`           // 余额
	UpdatedAt       time.Time      `json:"updated_at"`
}

// TableName 设置表名
func (ERC20TokenInfo) TableName() string {
	return "erc20_tokens"
}

func (ERC20Transaction) TableName() string {
	return "erc20_transactions"
}

func (ERC20Balance) TableName() string {
	return "erc20_balances"
}