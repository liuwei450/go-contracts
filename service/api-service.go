package service

import (
	"context"
	"database/sql"
	"fmt"
	"go-contracts/models"
	"go-contracts/synchronizer/node"
	"go-contracts/util"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type AirdropParams struct {
	// 具体参数根据实际业务定义
}

// GetBlockParams 获取区块信息的请求参数
type GetBlockParams struct {
	BlockNumber uint64 `json:"block_number"`
	BlockHash   string `json:"block_hash"`
}

// ERC20ContractParams ERC20合约相关操作的基础参数
type ERC20ContractParams struct {
	ContractAddress string `json:"contract_address"` // 合约地址
}

// ERC20AllowanceParams 查询授权额度的参数
type ERC20AllowanceParams struct {
	ERC20ContractParams
	Owner   string `json:"owner"`   // 授权方地址
	Spender string `json:"spender"` // 被授权方地址
}

// ERC20ApproveParams 设置授权的参数
type ERC20ApproveParams struct {
	ERC20ContractParams
	Spender string `json:"spender"` // 被授权方地址
	Value   string `json:"value"`   // 授权金额
}

// ERC20TransferParams 转账的参数
type ERC20TransferParams struct {
	ERC20ContractParams
	To    string `json:"to"`    // 接收方地址
	Value string `json:"value"` // 转账金额
}

// ERC20TransferFromParams 授权转账的参数
type ERC20TransferFromParams struct {
	ERC20ContractParams
	From  string `json:"from"`  // 发送方地址
	To    string `json:"to"`    // 接收方地址
	Value string `json:"value"` // 转账金额
}

// ERC20BalanceParams 查询余额的参数
type ERC20BalanceParams struct {
	ERC20ContractParams
	Account string `json:"account"` // 账户地址
}

type Service interface {
	AirdropBnb(ctx context.Context, params AirdropParams) error
	// 区块相关方法
	GetBlockByNumber(ctx context.Context, blockNumber uint64) (*models.Block, error)
	GetBlockByHash(ctx context.Context, blockHash string) (*models.Block, error)
	SaveBlock(ctx context.Context, block *models.Block) error
	GetLatestBlock(ctx context.Context) (*models.Block, error)

	// ERC20相关方法
	ERC20Allowance(ctx context.Context, params ERC20AllowanceParams) (*big.Int, error)
	ERC20Approve(ctx context.Context, params ERC20ApproveParams) (*big.Int, error)
	ERC20Transfer(ctx context.Context, params ERC20TransferParams) (*big.Int, error)
	ERC20TransferFrom(ctx context.Context, params ERC20TransferFromParams) (*big.Int, error)
	ERC20Balance(ctx context.Context, params ERC20BalanceParams) (*big.Int, error)
	ERC20TotalSupply(ctx context.Context, params ERC20ContractParams) (*big.Int, error)
	ERC20TokenInfo(ctx context.Context, params ERC20ContractParams) (*models.ERC20TokenInfo, error)
}

type serviceImpl struct {
	validator util.Validator
	db        *sql.DB
	// 区块链客户端接口
	ethClient node.EthClient
}

var _ Service = (*serviceImpl)(nil)

func New(validator util.Validator, ethClient node.EthClient) Service {
	return &serviceImpl{
		validator: validator,

		ethClient: ethClient,
	}
}
func (s *serviceImpl) AirdropBnb(ctx context.Context, params AirdropParams) error {
	// 业务逻辑
	// 使用 s.db 操作数据库等...

	return nil
}

// GetBlockByNumber 根据区块号获取区块信息
func (s *serviceImpl) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*models.Block, error) {
	var block models.Block
	// 实际实现：从数据库查询区块信息
	return &block, nil
}

// GetBlockByHash 根据区块哈希获取区块信息
func (s *serviceImpl) GetBlockByHash(ctx context.Context, blockHash string) (*models.Block, error) {
	var block models.Block
	// 实际实现：从数据库查询区块信息
	return &block, nil
}

// SaveBlock 保存区块信息到数据库
func (s *serviceImpl) SaveBlock(ctx context.Context, block *models.Block) error {
	// 实际实现：将区块信息保存到数据库
	return nil
}

// GetLatestBlock 获取最新区块信息
func (s *serviceImpl) GetLatestBlock(ctx context.Context) (*models.Block, error) {
	var block models.Block
	// 实际实现：获取数据库中最新的区块信息
	return &block, nil
}

// ERC20Allowance 查询授权额度
func (s *serviceImpl) ERC20Allowance(ctx context.Context, params ERC20AllowanceParams) (*big.Int, error) {
	contractAddress := common.HexToAddress(params.ContractAddress)
	owner := common.HexToAddress(params.Owner)
	spender := common.HexToAddress(params.Spender)

	return s.ethClient.ERC20Allowance(ctx, contractAddress, owner, spender)
}

// ERC20Approve 设置授权
func (s *serviceImpl) ERC20Approve(ctx context.Context, params ERC20ApproveParams) (*big.Int, error) {
	// 这里需要实现从某处获取TransactOpts的逻辑
	// 实际应用中通常会从配置或数据库中获取私钥
	// 为了示例，这里简化处理
	return nil, fmt.Errorf("未实现的方法: ERC20Approve")
}

// ERC20Transfer 转账
func (s *serviceImpl) ERC20Transfer(ctx context.Context, params ERC20TransferParams) (*big.Int, error) {
	// 这里需要实现从某处获取TransactOpts的逻辑
	// 实际应用中通常会从配置或数据库中获取私钥
	// 为了示例，这里简化处理
	return nil, fmt.Errorf("未实现的方法: ERC20Transfer")
}

// ERC20TransferFrom 授权转账
func (s *serviceImpl) ERC20TransferFrom(ctx context.Context, params ERC20TransferFromParams) (*big.Int, error) {
	// 这里需要实现从某处获取TransactOpts的逻辑
	// 实际应用中通常会从配置或数据库中获取私钥
	// 为了示例，这里简化处理
	return nil, fmt.Errorf("未实现的方法: ERC20TransferFrom")
}

// ERC20Balance 查询余额
func (s *serviceImpl) ERC20Balance(ctx context.Context, params ERC20BalanceParams) (*big.Int, error) {
	contractAddress := common.HexToAddress(params.ContractAddress)
	account := common.HexToAddress(params.Account)

	return s.ethClient.ERC20Balance(ctx, contractAddress, account)
}

// ERC20TotalSupply 查询总供应量
func (s *serviceImpl) ERC20TotalSupply(ctx context.Context, params ERC20ContractParams) (*big.Int, error) {
	contractAddress := common.HexToAddress(params.ContractAddress)

	return s.ethClient.ERC20TotalSupply(ctx, contractAddress)
}

// ERC20TokenInfo 获取代币信息
func (s *serviceImpl) ERC20TokenInfo(ctx context.Context, params ERC20ContractParams) (*models.ERC20TokenInfo, error) {
	contractAddress := common.HexToAddress(params.ContractAddress)

	name, symbol, decimals, err := s.ethClient.ERC20TokenInfo(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	totalSupply, err := s.ethClient.ERC20TotalSupply(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	return &models.ERC20TokenInfo{
			Address:     params.ContractAddress,
			Name:        name,
			Symbol:      symbol,
			Decimals:    decimals,
			TotalSupply: totalSupply.String(),
		},
		nil
}
