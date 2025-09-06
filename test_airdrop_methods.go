package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go-contracts/config"
	"go-contracts/contract"
	"go-contracts/database"
	"go-contracts/service"
	"go-contracts/synchronizer/node"
	"go-contracts/util"
	"math/big"
	"os"
	"os/signal"
	"syscall"
)

// 定义一个简单的模拟EthClient实现，用于测试
type mockEthClientImpl struct {}

func (m *mockEthClientImpl) GetERC20Contract(ctx context.Context, contractAddress common.Address) (*contract.Erc20, error) {
	return nil, nil
}

func (m *mockEthClientImpl) ERC20Allowance(ctx context.Context, contractAddress, owner, spender common.Address) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *mockEthClientImpl) ERC20Approve(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, spender common.Address, value *big.Int) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *mockEthClientImpl) ERC20Transfer(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, to common.Address, value *big.Int) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *mockEthClientImpl) ERC20TransferFrom(ctx context.Context, contractAddress common.Address, auth *bind.TransactOpts, from, to common.Address, value *big.Int) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *mockEthClientImpl) ERC20Balance(ctx context.Context, contractAddress common.Address, account common.Address) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *mockEthClientImpl) ERC20TotalSupply(ctx context.Context, contractAddress common.Address) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *mockEthClientImpl) ERC20TokenInfo(ctx context.Context, contractAddress common.Address) (string, string, uint8, error) {
	return "MockToken", "MOCK", 18, nil
}

func main() {
	// 初始化日志
	util.InitLogger()

	// 加载配置
	cfg, err := config.LoadConfig(nil)
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("接收到退出信号")
		cancel()
	}()
	defer cancel()

	// 初始化数据库
	db, err := database.NewDb(ctx, &cfg.MasterDB)
	if err != nil {
		fmt.Printf("初始化数据库失败: %v\n", err)
		return
	}
	defer db.Close()

	// 初始化Redis（如果需要）
	redisPool, err := database.NewRedis(nil, &cfg.Redis)
	if err != nil {
		fmt.Printf("初始化Redis失败: %v\n", err)
		return
	}
	defer redisPool.Close()

	// 初始化验证器
	var validator util.Validator

	// 初始化区块链客户端（使用我们定义的模拟客户端）
	var ethClient node.EthClient = &mockEthClientImpl{}

	// 创建服务实例
	svc := service.New(validator, ethClient)

	// 直接测试服务层的方法
	fmt.Println("===== 直接测试服务层方法 =====")

	// 测试AirdropGov方法
	govAddress, err := svc.AirdropGov(ctx)
	if err != nil {
		fmt.Printf("调用AirdropGov方法失败: %v\n", err)
	} else {
		fmt.Printf("AirdropGov方法返回: %s\n", govAddress)
	}

	// 测试AirdropSetGov方法
	params := service.AirdropSetGovParams{
		NewGov: "0x1234567890123456789012345678901234567890", // 示例地址
	}
	err = svc.AirdropSetGov(ctx, params)
	if err != nil {
		fmt.Printf("调用AirdropSetGov方法失败: %v\n", err)
	} else {
		fmt.Printf("AirdropSetGov方法调用成功\n")
	}

	fmt.Println("\n测试完成！")
}