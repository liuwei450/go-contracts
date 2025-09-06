package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
)

func main() {
	// 创建路由实例
	router := chi.NewRouter()

	// 注册与实际项目相同的路由
	healthPath := "/api/health"
	airdropSetGovPath := "/api/airdrop_set_gov"
	airdropGovPath := "/api/airdrop_gov"
	airdropBnbPath := "/api/airdrop_bnb"
	airdropErc20Path := "/api/airdrop_erc20"

	// 注册路由处理函数
	router.Get(healthPath, mockHandler)
	router.Post(airdropSetGovPath, mockHandler)
	router.Get(airdropGovPath, mockHandler)
	router.Get(airdropBnbPath, mockHandler)
	router.Post(airdropErc20Path, mockHandler)

	// 创建测试请求
	pathsToTest := []struct {
		path     string
		method   string
		expected bool
	}{
		{healthPath, "GET", true},
		{airdropSetGovPath, "POST", true},
		{airdropGovPath, "GET", true},
		{airdropBnbPath, "GET", true},
		{airdropErc20Path, "POST", true},
		{"/api/nonexistent", "GET", false},
	}

	fmt.Println("===== 测试路由匹配 =====")

	// 测试每个路径
	for _, test := range pathsToTest {
		// 创建测试请求
		req := httptest.NewRequest(test.method, test.path, nil)
		rec := httptest.NewRecorder()

		// 调用路由
		router.ServeHTTP(rec, req)

		// 检查响应
		statusCode := rec.Code
		matched := statusCode != http.StatusNotFound
		result := "失败"
		if matched == test.expected {
			result = "通过"
		}

		fmt.Printf("路径: %s, 方法: %s, 状态码: %d, 期望匹配: %v, 测试结果: %s\n",
			test.path, test.method, statusCode, test.expected, result)

		if matched {
			fmt.Printf("  响应内容: %s\n", rec.Body.String())
		}
	}

	fmt.Println("\n路由测试完成！")
}
