package service

import (
	"context"
	"fmt"
	"go-contracts/config"
	"go-contracts/contract"
	"go-contracts/models"
	"go-contracts/synchronizer/node"
	"go-contracts/util"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type AirdropParams struct {
	Recipients []string `json:"recipients"` // 接收者地址数组
	Amounts    []string `json:"amounts"`    // 金额数组（字符串形式）
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

// 设置空投合约授权地址的参数
type AirdropSetGovParams struct {
	NewGov string `json:"new_gov"` // 新的授权地址
}

type Service interface {
	AirdropBnb(ctx context.Context, params AirdropParams) error
	AirdropERC20(ctx context.Context, params AirdropParams) error
	AirdropSetGov(ctx context.Context, params AirdropSetGovParams) error
	AirdropGov(ctx context.Context) (string, error)
	// 区块相关方法
	GetBlockByNumber(ctx context.Context, blockNumber uint64) (*models.Block, error)
	GetBlockByHash(ctx context.Context, blockHash string) (*models.Block, error)
	SaveBlock(ctx context.Context, block *models.Block) error
	GetLatestBlock(ctx context.Context) (*models.Block, error)

	// ERC20相关方法
	ERC20Allowance(ctx context.Context, params ERC20AllowanceParams) (*big.Int, error)       // 查询授权额度
	ERC20Approve(ctx context.Context, params ERC20ApproveParams) (*big.Int, error)           // 授权
	ERC20Transfer(ctx context.Context, params ERC20TransferParams) (*big.Int, error)         // 转账
	ERC20TransferFrom(ctx context.Context, params ERC20TransferFromParams) (*big.Int, error) // 从授权地址转账
	ERC20Balance(ctx context.Context, params ERC20BalanceParams) (*big.Int, error)           // 查询余额
	ERC20TotalSupply(ctx context.Context, params ERC20ContractParams) (*big.Int, error)
	ERC20TokenInfo(ctx context.Context, params ERC20ContractParams) (*models.ERC20TokenInfo, error)
}

type serviceImpl struct {
	validator util.Validator

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
	// 1. 验证参数
	if len(params.Recipients) == 0 || len(params.Amounts) == 0 {
		return fmt.Errorf("接收者地址和金额不能为空")
	}
	if len(params.Recipients) != len(params.Amounts) {
		return fmt.Errorf("接收者地址数量和金额数量不匹配")
	}

	// 2. 连接到区块链节点
	client, err := ethclient.Dial(config.RAW_URL)
	if err != nil {
		return fmt.Errorf("连接区块链节点失败: %w", err)
	}
	defer client.Close()

	// 3. 解析私钥
	privateKey, err := crypto.HexToECDSA(config.PRIVATE_KEY)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %w", err)
	}

	// 4. 创建交易选项
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(97)) // 97是BSC测试网的链ID
	if err != nil {
		return fmt.Errorf("创建交易选项失败: %w", err)
	}

	// 5. 获取当前Gas价格
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		// 如果获取失败，使用默认值
		gasPrice = big.NewInt(10000000000) // 10 Gwei
		util.Log.Warn("获取Gas价格失败，使用默认值", "error", err)
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = 3000000 // 设置Gas上限

	// 6. 解析合约地址
	contractAddress := common.HexToAddress(config.AIRDROP_CONTRACT_ADDRESS)

	// 7. 创建合约实例
	airdropContract, err := contract.NewAirdropTransactor(contractAddress, client)
	if err != nil {
		return fmt.Errorf("创建空投合约实例失败: %w", err)
	}

	// 8. 转换接收者地址和金额
	recipients := make([]common.Address, len(params.Recipients))
	amounts := make([]*big.Int, len(params.Amounts))

	for i, addrStr := range params.Recipients {
		if !common.IsHexAddress(addrStr) {
			return fmt.Errorf("无效的以太坊地址: %s", addrStr)
		}
		recipients[i] = common.HexToAddress(addrStr)
	}

	for i, amountStr := range params.Amounts {
		amount, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			return fmt.Errorf("无效的金额格式: %s", amountStr)
		}
		amounts[i] = amount
	}

	// 9. 计算总金额并设置交易价值
	totalAmount := new(big.Int)
	for _, amount := range amounts {
		totalAmount.Add(totalAmount, amount)
	}
	auth.Value = totalAmount

	// 10. 调用空投合约方法
	tx, err := airdropContract.AirdropBNB(auth, recipients, amounts)
	if err != nil {
		return fmt.Errorf("调用空投合约失败: %w", err)
	}

	// 11. 记录交易信息
	util.Log.Info("BNB空投交易已发送", "txHash", tx.Hash().Hex())

	// 12. 可以在这里添加数据库记录逻辑
	// 例如保存交易哈希、接收者、金额等信息到数据库

	return nil
}

// AirdropERC20 实现ERC20代币空投功能
func (s *serviceImpl) AirdropERC20(ctx context.Context, params AirdropParams) error {
	// 1. 验证参数
	if len(params.Recipients) == 0 || len(params.Amounts) == 0 {
		return fmt.Errorf("接收者地址和金额不能为空")
	}
	if len(params.Recipients) != len(params.Amounts) {
		return fmt.Errorf("接收者地址数量和金额数量不匹配")
	}

	// 2. 连接到区块链节点
	client, err := ethclient.Dial(config.RAW_URL)
	if err != nil {
		return fmt.Errorf("连接区块链节点失败: %w", err)
	}
	defer client.Close()

	// 3. 解析私钥
	privateKey, err := crypto.HexToECDSA(config.PRIVATE_KEY)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %w", err)
	}

	// 4. 创建交易选项
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(97)) // 97是BSC测试网的链ID
	if err != nil {
		return fmt.Errorf("创建交易选项失败: %w", err)
	}

	// 5. 获取当前Gas价格
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		// 如果获取失败，使用默认值
		gasPrice = big.NewInt(10000000000) // 10 Gwei
		util.Log.Warn("获取Gas价格失败，使用默认值", "error", err)
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = 3000000 // 设置Gas上限

	// 6. 解析合约地址
	contractAddress := common.HexToAddress(config.AIRDROP_CONTRACT_ADDRESS)

	// 7. 创建合约实例
	airdropContract, err := contract.NewAirdropTransactor(contractAddress, client)
	if err != nil {
		return fmt.Errorf("创建空投合约实例失败: %w", err)
	}

	// 8. 转换接收者地址和金额
	recipients := make([]common.Address, len(params.Recipients))
	amounts := make([]*big.Int, len(params.Amounts))

	for i, addrStr := range params.Recipients {
		if !common.IsHexAddress(addrStr) {
			return fmt.Errorf("无效的以太坊地址: %s", addrStr)
		}
		recipients[i] = common.HexToAddress(addrStr)
	}

	for i, amountStr := range params.Amounts {
		amount, ok := new(big.Int).SetString(amountStr, 10)
		if !ok {
			return fmt.Errorf("无效的金额格式: %s", amountStr)
		}
		amounts[i] = amount
	}

	// 9. 调用ERC20空投合约方法
	tx, err := airdropContract.AirdropERC20(auth, recipients, amounts)
	if err != nil {
		return fmt.Errorf("调用ERC20空投合约失败: %w", err)
	}

	// 10. 记录交易信息
	util.Log.Info("ERC20空投交易已发送", "txHash", tx.Hash().Hex())

	// 11. 可以在这里添加数据库记录逻辑
	// 例如保存交易哈希、接收者、金额等信息到数据库

	return nil
}

// GetBlockByNumber 根据区块号获取区块信息
func (s *serviceImpl) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*models.Block, error) {
	// 检查 service 是否为 nil
	if s == nil {
		return nil, fmt.Errorf("service is nil")
	}

	// 检查区块链客户端是否存在
	if s.ethClient == nil {
		return nil, fmt.Errorf("区块链客户端未初始化")
	}

	// 注意：当前的 EthClient 接口不直接支持获取区块信息
	// 这里仅返回一个示例区块数据用于演示
	now := time.Now()
	result := &models.Block{
		BlockNumber: blockNumber,
		BlockHash:   fmt.Sprintf("0xblockhash%v", blockNumber),
		ParentHash:  fmt.Sprintf("0xparenthash%v", blockNumber-1),
		Time:        uint64(now.Unix()),
		Timestamp:   now,
	}

	return result, nil
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

// AirdropSetGov 设置空投合约授权地址
func (s *serviceImpl) AirdropSetGov(ctx context.Context, params AirdropSetGovParams) error {
	// 1. 连接到区块链节点
	client, err := ethclient.Dial(config.RAW_URL)
	if err != nil {
		return fmt.Errorf("连接区块链节点失败: %w", err)
	}
	defer client.Close()

	// 2. 解析私钥
	privateKey, err := crypto.HexToECDSA(config.PRIVATE_KEY)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %w", err)
	}

	// 3. 创建交易选项
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(97)) // 97是BSC测试网的链ID
	if err != nil {
		return fmt.Errorf("创建交易选项失败: %w", err)
	}

	// 4. 获取当前Gas价格
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		// 如果获取失败，使用默认值
		gasPrice = big.NewInt(10000000000) // 10 Gwei
		util.Log.Warn("获取Gas价格失败，使用默认值", "error", err)
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = 3000000 // 设置Gas上限

	// 5. 解析合约地址
	contractAddress := common.HexToAddress(config.AIRDROP_CONTRACT_ADDRESS)

	// 6. 创建合约实例
	airdropContract, err := contract.NewAirdropTransactor(contractAddress, client)
	if err != nil {
		return fmt.Errorf("创建空投合约实例失败: %w", err)
	}

	// 7. 验证并解析新的授权地址
	if !common.IsHexAddress(params.NewGov) {
		return fmt.Errorf("无效的以太坊地址: %s", params.NewGov)
	}
	newGovAddr := common.HexToAddress(params.NewGov)

	// 8. 调用setGov方法
	tx, err := airdropContract.SetGov(auth, newGovAddr)
	if err != nil {
		return fmt.Errorf("调用setGov方法失败: %w", err)
	}

	// 9. 记录交易信息
	util.Log.Info("设置空投合约授权地址交易已发送", "txHash", tx.Hash().Hex(), "newGov", params.NewGov)

	return nil
}

// AirdropGov 查询空投合约授权地址
func (s *serviceImpl) AirdropGov(ctx context.Context) (string, error) {
	// 1. 连接到区块链节点
	client, err := ethclient.Dial(config.RAW_URL)
	if err != nil {
		return "", fmt.Errorf("连接区块链节点失败: %w", err)
	}
	defer client.Close()

	// 2. 解析合约地址
	contractAddress := common.HexToAddress(config.AIRDROP_CONTRACT_ADDRESS)

	// 3. 创建合约实例（只读）
	airdropContract, err := contract.NewAirdropCaller(contractAddress, client)
	if err != nil {
		return "", fmt.Errorf("创建空投合约只读实例失败: %w", err)
	}

	// 4. 调用gov方法查询授权地址
	govAddr, err := airdropContract.Gov(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", fmt.Errorf("查询授权地址失败: %w", err)
	}

	// 5. 返回授权地址
	return govAddr.Hex(), nil
}
