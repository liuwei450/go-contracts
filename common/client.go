package common

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ethClient  *ethclient.Client
	clientOnce sync.Once
)

// InitEthClient 初始化以太坊/BSC客户端连接
func InitEthClient(rpcURL string) error {
	var err error
	clientOnce.Do(func() {
		// 连接到BSC网络节点
		ethClient, err = ethclient.Dial(rpcURL)
		if err != nil {
			err = fmt.Errorf("连接BSC节点失败: %w", err)
			return
		}
	})
	return err
}

// GetEthClient 获取以太坊/BSC客户端实例
func GetEthClient() *ethclient.Client {
	return ethClient
}

// CloseEthClient 关闭以太坊/BSC客户端连接
func CloseEthClient() error {
	if ethClient != nil {
		// 在Go-ethereum v1.16.2中，ethclient.Client.Close()不返回错误
		ethClient.Close()
		ethClient = nil
	}
	return nil
}
