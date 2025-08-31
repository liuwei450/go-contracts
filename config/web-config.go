package config

import "go-public-web/common"

const (
	// 部署的链的 RPC
	RAW_URL                  = "https://data-seed-prebsc-2-s1.binance.org:8545"
	PRIVATE_KEY              = common.PRIVATE_KEY
	ERC20_CONTRACT_ADDRESS   = common.ERC20_CONTRACT_ADDRESS
	AIRDROP_CONTRACT_ADDRESS = common.AIRDROP_CONTRACT_ADDRESS
	MTK_CONTRACT_ADDRESS     = common.MTK_CONTRACT_ADDRESS
	// 原生代币地址（全零地址，用于表示BNB等原生代币）
	NATIVE_TOKEN_ADDRESS = "0x0000000000000000000000000000000000000000"
)
