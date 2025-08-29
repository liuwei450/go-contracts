package node

import (
	"context"
	"github.com/urfave/cli/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go-contracts/contract"
	"math/big"
)

// EthClient 以太坊/BSC客户端接口
type EthClient interface {
	// 获取ERC20合约实例
	GetERC20Contract(ctx context.Context, contractAddress common.Address) (*contract.Erc20, error)
	// 查询授权额度
	ERC20Allowance(ctx context.Context, contractAddress common.Address, owner, spender common.Address) (*big.Int, error)
	// 设置授权
	ERC20Approve(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, spender common.Address, value *big.Int) (*big.Int, error)
	// 转账
	ERC20Transfer(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, to common.Address, value *big.Int) (*big.Int, error)
	// 授权转账
	ERC20TransferFrom(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, from, to common.Address, value *big.Int) (*big.Int, error)
	// 查询余额
	ERC20Balance(ctx context.Context, contractAddress common.Address, account common.Address) (*big.Int, error)
	// 查询总供应量
	ERC20TotalSupply(ctx context.Context, contractAddress common.Address) (*big.Int, error)
	// 获取代币信息
	ERC20TokenInfo(ctx context.Context, contractAddress common.Address) (string, string, uint8, error)
}

// ethClientImpl 实现EthClient接口
type ethClientImpl struct {
	client *ethclient.Client
}

// NewEthClientImpl 创建新的EthClient实现
func NewEthClientImpl(client *ethclient.Client) EthClient {
	return &ethClientImpl{
		client: client,
	}
}

// GetERC20Contract 获取ERC20合约实例
func (e *ethClientImpl) GetERC20Contract(ctx context.Context, contractAddress common.Address) (*contract.Erc20, error) {
	return contract.NewErc20(contractAddress, e.client)
}

// ERC20Allowance 查询授权额度
func (e *ethClientImpl) ERC20Allowance(ctx context.Context, contractAddress common.Address, owner, spender common.Address) (*big.Int, error) {
	contract, err := e.GetERC20Contract(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	return contract.Allowance(&bind.CallOpts{Context: ctx}, owner, spender)
}

// ERC20Approve 设置授权
func (e *ethClientImpl) ERC20Approve(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, spender common.Address, value *big.Int) (*big.Int, error) {
	contract, err := e.GetERC20Contract(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	tx, err := contract.Approve(auth, spender, value)
	if err != nil {
		return nil, err
	}

	// 返回交易哈希的字节表示作为ID
	return new(big.Int).SetBytes(tx.Hash().Bytes()), nil
}

// ERC20Transfer 转账
func (e *ethClientImpl) ERC20Transfer(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, to common.Address, value *big.Int) (*big.Int, error) {
	contract, err := e.GetERC20Contract(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	tx, err := contract.Transfer(auth, to, value)
	if err != nil {
		return nil, err
	}

	// 返回交易哈希的字节表示作为ID
	return new(big.Int).SetBytes(tx.Hash().Bytes()), nil
}

// ERC20TransferFrom 授权转账
func (e *ethClientImpl) ERC20TransferFrom(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, from, to common.Address, value *big.Int) (*big.Int, error) {
	contract, err := e.GetERC20Contract(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	tx, err := contract.TransferFrom(auth, from, to, value)
	if err != nil {
		return nil, err
	}

	// 返回交易哈希的字节表示作为ID
	return new(big.Int).SetBytes(tx.Hash().Bytes()), nil
}

// ERC20Balance 查询余额
func (e *ethClientImpl) ERC20Balance(ctx context.Context, contractAddress common.Address, account common.Address) (*big.Int, error) {
	contract, err := e.GetERC20Contract(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	return contract.BalanceOf(&bind.CallOpts{Context: ctx}, account)
}

// ERC20TotalSupply 查询总供应量
func (e *ethClientImpl) ERC20TotalSupply(ctx context.Context, contractAddress common.Address) (*big.Int, error) {
	contract, err := e.GetERC20Contract(ctx, contractAddress)
	if err != nil {
		return nil, err
	}

	return contract.TotalSupply(&bind.CallOpts{Context: ctx})
}

// ERC20TokenInfo 获取代币信息
func (e *ethClientImpl) ERC20TokenInfo(ctx context.Context, contractAddress common.Address) (string, string, uint8, error) {
	contract, err := e.GetERC20Contract(ctx, contractAddress)
	if err != nil {
		return "", "", 0, err
	}

	name, err := contract.Name(&bind.CallOpts{Context: ctx})
	if err != nil {
		return "", "", 0, err
	}

	symbol, err := contract.Symbol(&bind.CallOpts{Context: ctx})
	if err != nil {
		return name, "", 0, err
	}

	decimals, err := contract.Decimals(&bind.CallOpts{Context: ctx})
	if err != nil {
		return name, symbol, 0, err
	}

	return name, symbol, decimals, nil
}

// DialEthClient 连接以太坊/BSC节点
func DialEthClient(ctx cli.Context, rpcUrl string) (EthClient, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, err
	}

	return NewEthClientImpl(client), nil
}
