package worker

import (
	"context"
	"fmt"
	"go-contracts/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

// ERC20Worker 提供ERC20代币相关的操作功能
type ERC20Worker struct {
	client *ethclient.Client
}

// NewERC20Worker 创建新的ERC20Worker实例
func NewERC20Worker(client *ethclient.Client) *ERC20Worker {
	return &ERC20Worker{
		client: client,
	}
}

// GetContract 获取ERC20合约实例
func (w *ERC20Worker) GetContract(contractAddress common.Address) (*contract.Erc20, error) {
	return contract.NewErc20(contractAddress, w.client)
}

// GetAllowance 查询授权额度
func (w *ERC20Worker) GetAllowance(ctx context.Context, contractAddress common.Address, owner, spender common.Address) (*big.Int, error) {
	contract, err := w.GetContract(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("获取合约实例失败: %w", err)
	}

	return contract.Allowance(&bind.CallOpts{Context: ctx}, owner, spender)
}

// Approve 设置授权
func (w *ERC20Worker) Approve(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, spender common.Address, value *big.Int) (*big.Int, error) {
	contract, err := w.GetContract(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("获取合约实例失败: %w", err)
	}

	tx, err := contract.Approve(auth, spender, value)
	if err != nil {
		return nil, fmt.Errorf("授权失败: %w", err)
	}

	// 返回交易哈希的字节表示作为ID
	return new(big.Int).SetBytes(tx.Hash().Bytes()), nil
}

// Transfer 转账
func (w *ERC20Worker) Transfer(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, to common.Address, value *big.Int) (*big.Int, error) {
	contract, err := w.GetContract(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("获取合约实例失败: %w", err)
	}

	tx, err := contract.Transfer(auth, to, value)
	if err != nil {
		return nil, fmt.Errorf("转账失败: %w", err)
	}

	// 返回交易哈希的字节表示作为ID
	return new(big.Int).SetBytes(tx.Hash().Bytes()), nil
}

// TransferFrom 授权转账
func (w *ERC20Worker) TransferFrom(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, from, to common.Address, value *big.Int) (*big.Int, error) {
	contract, err := w.GetContract(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("获取合约实例失败: %w", err)
	}

	tx, err := contract.TransferFrom(auth, from, to, value)
	if err != nil {
		return nil, fmt.Errorf("授权转账失败: %w", err)
	}

	// 返回交易哈希的字节表示作为ID
	return new(big.Int).SetBytes(tx.Hash().Bytes()), nil
}

// GetBalance 查询余额
func (w *ERC20Worker) GetBalance(ctx context.Context, contractAddress common.Address, account common.Address) (*big.Int, error) {
	contract, err := w.GetContract(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("获取合约实例失败: %w", err)
	}

	return contract.BalanceOf(&bind.CallOpts{Context: ctx}, account)
}

// GetTotalSupply 查询总供应量
func (w *ERC20Worker) GetTotalSupply(ctx context.Context, contractAddress common.Address) (*big.Int, error) {
	contract, err := w.GetContract(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("获取合约实例失败: %w", err)
	}

	return contract.TotalSupply(&bind.CallOpts{Context: ctx})
}

// GetTokenInfo 获取代币信息
func (w *ERC20Worker) GetTokenInfo(ctx context.Context, contractAddress common.Address) (string, string, uint8, error) {
	contract, err := w.GetContract(contractAddress)
	if err != nil {
		return "", "", 0, fmt.Errorf("获取合约实例失败: %w", err)
	}

	name, err := contract.Name(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", "", 0, fmt.Errorf("获取代币名称失败: %w", err)
	}

	symbol, err := contract.Symbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return name, "", 0, fmt.Errorf("获取代币符号失败: %w", err)
	}

	decimals, err := contract.Decimals(&bind.CallOpts{Context: ctx})
	if err != nil {
		return name, symbol, 0, fmt.Errorf("获取代币小数位失败: %w", err)
	}

	return name, symbol, decimals, nil
}