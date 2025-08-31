package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

// Validator 提供各种验证功能的结构
// 用于验证区块链地址、交易哈希、数字范围等（适用于币安智能链BSC）
// 目前版本: v1.0.0
// 最后更新: 2025-08-21
// 作者: Trae AI
// 许可证: MIT
// 依赖: github.com/ethereum/go-ethereum
// 使用示例: validator := util.NewValidator()
//          isValid := validator.IsValidAddress("0xa8aa61bf1c35eceb56d9bffb2f59ad34898a1dbb")
type Validator struct{}

// NewValidator 创建一个新的验证器实例
// 返回: *Validator - 验证器指针
func NewValidator() *Validator {
	return &Validator{}
}

// IsValidAddress 验证区块链地址格式是否正确（适用于币安智能链BSC）
// 参数: address string - 要验证的区块链地址
// 返回: bool - 地址是否有效
func (v *Validator) IsValidAddress(address string) bool {
	// 检查地址长度和前缀
	if !strings.HasPrefix(address, "0x") || len(address) != 42 {
		return false
	}

	// 检查地址是否只包含有效的十六进制字符
	hexRegex := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return hexRegex.MatchString(address)
}

// IsValidTransactionHash 验证区块链交易哈希格式是否正确（适用于币安智能链BSC）
// 参数: hash string - 要验证的交易哈希
// 返回: bool - 哈希是否有效
func (v *Validator) IsValidTransactionHash(hash string) bool {
	// 检查哈希长度和前缀
	if !strings.HasPrefix(hash, "0x") || len(hash) != 66 {
		return false
	}

	// 检查哈希是否只包含有效的十六进制字符
	hexRegex := regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	return hexRegex.MatchString(hash)

}

// IsValidBlockHash 验证区块链区块哈希格式是否正确（适用于币安智能链BSC）
// 参数: hash string - 要验证的区块哈希
// 返回: bool - 哈希是否有效
func (v *Validator) IsValidBlockHash(hash string) bool {
	// 区块哈希验证与交易哈希相同
	return v.IsValidTransactionHash(hash)
}

// IsValidBlockNumber 验证区块号是否为有效的正整数
// 参数: blockNumber string - 要验证的区块号
// 返回: bool - 区块号是否有效
func (v *Validator) IsValidBlockNumber(blockNumber string) bool {
	if blockNumber == "latest" || blockNumber == "earliest" || blockNumber == "pending" {
		return true
	}

	num, err := strconv.ParseUint(blockNumber, 10, 64)
	return err == nil && num >= 0
}

// IsValidPrivateKey 验证区块链私钥格式是否正确（适用于币安智能链BSC）
// 参数: privateKey string - 要验证的私钥
// 返回: bool - 私钥是否有效, error - 错误信息
func (v *Validator) IsValidPrivateKey(privateKey string) (bool, error) {
	// 尝试解析私钥
	_, err := crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
	if err != nil {
		return false, errors.New("invalid private key format")
	}
	return true, nil
}

// IsValidPublicKey 验证区块链公钥格式是否正确（适用于币安智能链BSC）
// 参数: publicKey string - 要验证的公钥
// 返回: bool - 公钥是否有效, error - 错误信息
func (v *Validator) IsValidPublicKey(publicKey string) (bool, error) {
	// 移除前缀(如果有)
	publicKey = strings.TrimPrefix(publicKey, "0x")
	// 公钥应该是130个字符(不包括0x前缀)
	if len(publicKey) != 130 {
		return false, errors.New("invalid public key length")
	}

	// 检查公钥是否只包含有效的十六进制字符
	hexRegex := regexp.MustCompile("^[0-9a-fA-F]{130}$")
	if !hexRegex.MatchString(publicKey) {
		return false, errors.New("invalid public key format")
	}

	// 检查公钥是否以04开头(非压缩公钥)
	if !strings.HasPrefix(strings.ToLower(publicKey), "04") {
		return false, errors.New("public key should start with 04")
	}

	return true, nil
}

// IsValidAmount 验证金额是否为有效的正数
// 参数: amount string - 要验证的金额
// 返回: bool - 金额是否有效
func (v *Validator) IsValidAmount(amount string) bool {
	// 允许0或正数
	amountRegex := regexp.MustCompile(`^\d+(\.\d+)?$`)
	return amountRegex.MatchString(amount)
}

// IsInRange 验证数值是否在指定范围内
// 参数: value int64 - 要验证的数值
// 参数: min int64 - 最小值
// 参数: max int64 - 最大值
// 返回: bool - 数值是否在范围内
func (v *Validator) IsInRange(value, min, max int64) bool {
	return value >= min && value <= max
}

// IsValidERC20Symbol 验证ERC20代币符号格式是否正确
// 参数: symbol string - 要验证的代币符号
// 返回: bool - 符号是否有效
func (v *Validator) IsValidERC20Symbol(symbol string) bool {
	// 代币符号应该是1-10个大写字母
	symbolRegex := regexp.MustCompile("^[A-Z]{1,10}$")
	return symbolRegex.MatchString(symbol)
}

// IsValidHex 验证字符串是否为有效的十六进制格式
// 参数: hex string - 要验证的十六进制字符串
// 参数: hasPrefix bool - 是否需要0x前缀
// 返回: bool - 字符串是否为有效的十六进制格式
func (v *Validator) IsValidHex(hex string, hasPrefix bool) bool {
	if hasPrefix {
		if !strings.HasPrefix(hex, "0x") {
			return false
		}
		hex = hex[2:]
	}

	// 检查十六进制字符
	hexRegex := regexp.MustCompile("^[0-9a-fA-F]+$")
	return hexRegex.MatchString(hex)
}
