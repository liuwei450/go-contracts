package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	baseURL := "http://localhost:8090"

	// 测试健康检查接口
	testHealthCheck(baseURL)

	// 测试ERC20余额查询接口
	testERC20Balance(baseURL)

	// 测试ERC20总供应量接口
	testERC20TotalSupply(baseURL)

	// 测试ERC20代币信息接口
	testERC20TokenInfo(baseURL)

	fmt.Println("所有接口测试完成！")
}

func testHealthCheck(baseURL string) {
	url := baseURL + "/api/health"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("健康检查接口测试失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("健康检查接口: StatusCode=%d, Response=%s\n", resp.StatusCode, body)
}

func testERC20Balance(baseURL string) {
	url := baseURL + "/api/erc20/balance"
	params := map[string]string{
		"contract_address": "0x0000000000000000000000000000000000000000", // 测试地址
		"account":          "0x1234567890123456789012345678901234567890",
	}
	testPostAPI(url, params, "ERC20余额查询接口")
}

func testERC20TotalSupply(baseURL string) {
	url := baseURL + "/api/erc20/total_supply"
	params := map[string]string{
		"contract_address": "0x0000000000000000000000000000000000000000", // 测试地址
	}
	testPostAPI(url, params, "ERC20总供应量接口")
}

func testERC20TokenInfo(baseURL string) {
	url := baseURL + "/api/erc20/token_info"
	params := map[string]string{
		"contract_address": "0x0000000000000000000000000000000000000000", // 测试地址
	}
	testPostAPI(url, params, "ERC20代币信息接口")
}

func testPostAPI(url string, params map[string]string, apiName string) {
	jsonData, _ := json.Marshal(params)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("%s测试失败: %v\n", apiName, err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s: StatusCode=%d, Response=%s\n", apiName, resp.StatusCode, body)
}