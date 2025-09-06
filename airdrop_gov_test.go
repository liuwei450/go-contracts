package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	// 从命令行参数获取API地址，如果没有提供则使用默认地址
	baseURL := "http://localhost:8090"
	if len(os.Args) > 1 {
		baseURL = os.Args[1]
	}

	fmt.Println("===== 空投合约授权地址接口测试工具 =====")

	// 1. 先查询当前的授权地址
	testAirdropGov(baseURL)

	// 2. 提示用户可以选择设置新的授权地址
	fmt.Println("\n注意：由于设置新的授权地址需要真实的区块链交易，如果要测试设置功能，请先确保环境配置正确（如私钥、RPC节点等）。")
	fmt.Println("使用方式：")
	fmt.Println("1. 查询授权地址: go run test_airdrop_gov.go [API地址]")
	fmt.Println("2. 设置授权地址: 修改此脚本中的注释代码以启用设置功能")

	fmt.Println("\n测试完成！")
}

// 测试查询空投合约授权地址接口
func testAirdropGov(baseURL string) {
	url := baseURL + "/api/airdrop_gov"
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("查询空投合约授权地址接口测试失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("查询空投合约授权地址接口: StatusCode=%d, Response=%s\n", resp.StatusCode, body)
}

// 测试设置空投合约授权地址接口（需要真实环境才能完全测试）
func testAirdropSetGov(baseURL string, newGovAddress string) {
	url := baseURL + "/api/airdrop_set_gov"
	params := map[string]string{
		"new_gov": newGovAddress,
	}
	jsonData, _ := json.Marshal(params)
	client := &http.Client{Timeout: 30 * time.Second} // 设置更长的超时时间，因为区块链交易可能需要时间
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("设置空投合约授权地址接口测试失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("设置空投合约授权地址接口: StatusCode=%d, Response=%s\n", resp.StatusCode, body)

	// 设置后再次查询，验证是否设置成功
	fmt.Println("\n设置后再次查询授权地址：")
	testAirdropGov(baseURL)
}