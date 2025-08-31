package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidator_IsValidAddress 测试区块链地址验证功能（适用于币安智能链BSC）
func TestValidator_IsValidAddress(t *testing.T) {
	validator := NewValidator()

	// 测试用例
	testCases := []struct {
		name     string
		address  string
		expected bool
	}{{
		name:     "有效的BSC地址",
		address:  "0xa8aa61bf1c35eceb56d9bffb2f59ad34898a1dbb",
		expected: true,
	}, {
		name:     "有效的以太坊地址格式",
		address:  "0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
		expected: true,
	}, {
		name:     "无效的地址 - 缺少0x前缀",
		address:  "71C7656EC7ab88b098defB751B7401B5f6d8976F",
		expected: false,
	}, {
		name:     "无效的地址 - 长度错误",
		address:  "0x71C7656EC7ab88b098defB751B7401B5f6d897",
		expected: false,
	}, {
		name:     "无效的地址 - 包含非法字符",
		address:  "0x71C7656EC7ab88b098defB751B7401B5f6d8976G",
		expected: false,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.IsValidAddress(tc.address)
			assert.Equal(t, tc.expected, result, tc.name)
		})
	}
}

// TestValidator_IsValidTransactionHash 测试区块链交易哈希验证功能（适用于币安智能链BSC）
func TestValidator_IsValidTransactionHash(t *testing.T) {
	validator := NewValidator()

	// 测试用例
	testCases := []struct {
		name     string
		hash     string
		expected bool
	}{{
		name:     "有效的交易哈希",
		hash:     "0x0a8a3f149a187c5a2648277181aa91569130165f1876578751f3f595d1e05fc9",
		expected: true,
	}, {
		name:     "无效的哈希 - 缺少0x前缀",
		hash:     "0a8a3f149a187c5a2648277181aa91569130165f1876578751f3f595d1e05fc9",
		expected: false,
	}, {
		name:     "无效的哈希 - 长度错误",
		hash:     "0x0a8a3f149a187c5a2648277181aa91569130165f1876578751f3f595d1e05fc",
		expected: false,
	}, {
		name:     "无效的哈希 - 包含非法字符",
		hash:     "0x0a8a3f149a187c5a2648277181aa91569130165f1876578751f3f595d1e05fcG",
		expected: false,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.IsValidTransactionHash(tc.hash)
			assert.Equal(t, tc.expected, result, tc.name)
		})
	}
}

// TestValidator_IsValidBlockNumber 测试区块链区块号验证功能（适用于币安智能链BSC）
func TestValidator_IsValidBlockNumber(t *testing.T) {
	validator := NewValidator()

	// 测试用例
	testCases := []struct {
		name     string
		number   string
		expected bool
	}{{
		name:     "有效的区块号 - 数字",
		number:   "12345678",
		expected: true,
	}, {
		name:     "有效的区块号 - latest",
		number:   "latest",
		expected: true,
	}, {
		name:     "有效的区块号 - earliest",
		number:   "earliest",
		expected: true,
	}, {
		name:     "有效的区块号 - pending",
		number:   "pending",
		expected: true,
	}, {
		name:     "无效的区块号 - 负数",
		number:   "-12345",
		expected: false,
	}, {
		name:     "无效的区块号 - 非数字",
		number:   "abcdef",
		expected: false,
	}, {
		name:     "无效的区块号 - 空字符串",
		number:   "",
		expected: false,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.IsValidBlockNumber(tc.number)
			assert.Equal(t, tc.expected, result, tc.name)
		})
	}
}

// TestValidator_IsValidAmount 测试金额验证功能
func TestValidator_IsValidAmount(t *testing.T) {
	validator := NewValidator()

	// 测试用例
	testCases := []struct {
		name     string
		amount   string
		expected bool
	}{{
		name:     "有效的金额 - 整数",
		amount:   "12345",
		expected: true,
	}, {
		name:     "有效的金额 - 小数",
		amount:   "12345.6789",
		expected: true,
	}, {
		name:     "有效的金额 - 零",
		amount:   "0",
		expected: true,
	}, {
		name:     "无效的金额 - 负数",
		amount:   "-12345",
		expected: false,
	}, {
		name:     "无效的金额 - 非数字",
		amount:   "abcdef",
		expected: false,
	}, {
		name:     "无效的金额 - 多个小数点",
		amount:   "123.45.67",
		expected: false,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.IsValidAmount(tc.amount)
			assert.Equal(t, tc.expected, result, tc.name)
		})
	}
}

// TestValidator_IsValidHex 测试十六进制字符串验证功能
func TestValidator_IsValidHex(t *testing.T) {
	validator := NewValidator()

	// 测试用例
	testCases := []struct {
			name      string
			hasHex    string
			hasPrefix bool
			expected  bool
		}{{
			name:      "有效的十六进制 - 带前缀",
			hasHex:    "0x123ABC",
			hasPrefix: true,
			expected:  true,
		}, {
			name:      "有效的十六进制 - 不带前缀",
			hasHex:    "123ABC",
			hasPrefix: false,
			expected:  true,
		}, {
			name:      "无效的十六进制 - 带前缀但格式错误",
			hasHex:    "0x123ABCZ",
			hasPrefix: true,
			expected:  false,
		}, {
			name:      "无效的十六进制 - 不带前缀但需要前缀",
			hasHex:    "123ABC",
			hasPrefix: true,
			expected:  false,
		}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.IsValidHex(tc.hasHex, tc.hasPrefix)
			assert.Equal(t, tc.expected, result, tc.name)
		})
	}
}