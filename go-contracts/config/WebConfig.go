package config

const (
	// 部署的链的 RPC
	RPC_URL = "https://bsc-testnet.publicnode.com"

	//账户的私钥（十六进制格式），用于签名交易
	PRIVATE_KEY_HEX = "0xa8Aa61bF1c35ECEb56d9Bffb2f59AD34898A1dbb"
	//部署好的质押合约地址
	CONTRACT_ADRESS = "0x81C0C4d1a073A77Eb3f6C4a81c4aC04dbf35E8dC"
	ERC20_ADDRESS   = "0xfea7f4509F643038B5A9C1220BF0f5c661FA7ffc" //  ERC20 地址
	CHAIN_ID        = 97                                           // BSC 测试网
	MYSQL_DSN       = "root:123456@tcp(127.0.0.1:3306)/user_db?charset=utf8mb4&parseTime=True&loc=Local"
	ABI_JSON        = `[
			{
				"inputs": [
					{
						"internalType": "contract IERC20",
						"name": "_mtkToken",
						"type": "address"
					}
				],
				"stateMutability": "nonpayable",
				"type": "constructor"
			},
			{
				"inputs": [],
				"name": "StakeNotFound",
				"type": "error"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "user",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "stakeId",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "enum MtkContracts.StakingPeriod",
						"name": "period",
						"type": "uint8"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "timestamp",
						"type": "uint256"
					}
				],
				"name": "Staked",
				"type": "event"
			},
			{
				"anonymous": false,
				"inputs": [
					{
						"indexed": true,
						"internalType": "address",
						"name": "user",
						"type": "address"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "stakeId",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "principal",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "totalAmount",
						"type": "uint256"
					},
					{
						"indexed": false,
						"internalType": "uint256",
						"name": "stakeIndex",
						"type": "uint256"
					}
				],
				"name": "Withdrawn",
				"type": "event"
			},
			{
				"inputs": [
					{
						"internalType": "enum MtkContracts.StakingPeriod",
						"name": "",
						"type": "uint8"
					}
				],
				"name": "apy",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "enum MtkContracts.StakingPeriod",
						"name": "",
						"type": "uint8"
					}
				],
				"name": "durations",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					},
					{
						"internalType": "enum MtkContracts.StakingPeriod",
						"name": "period",
						"type": "uint8"
					}
				],
				"name": "stake",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"name": "stakeIdToOwner",
				"outputs": [
					{
						"internalType": "address",
						"name": "",
						"type": "address"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [],
				"name": "stakingToken",
				"outputs": [
					{
						"internalType": "contract IERC20",
						"name": "",
						"type": "address"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "address",
						"name": "",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "",
						"type": "uint256"
					}
				],
				"name": "userStakes",
				"outputs": [
					{
						"internalType": "uint256",
						"name": "stakeId",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "amount",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "startTime",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "endTime",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "rewardRate",
						"type": "uint256"
					},
					{
						"internalType": "bool",
						"name": "isActive",
						"type": "bool"
					},
					{
						"internalType": "enum MtkContracts.StakingPeriod",
						"name": "period",
						"type": "uint8"
					}
				],
				"stateMutability": "view",
				"type": "function"
			},
			{
				"inputs": [
					{
						"internalType": "uint256",
						"name": "stakeId",
						"type": "uint256"
					}
				],
				"name": "withdraw",
				"outputs": [],
				"stateMutability": "nonpayable",
				"type": "function"
			}
		]`
)
