package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-contracts/config"
	"go-contracts/router"
	"go-contracts/service"
	"go-contracts/models"
	"net/http"
	"net/http/httptest"
	"reflect"
	"math/big"
)

// Mock handler for testing
func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Handler called: " + r.URL.Path))
}

// MockService 实现service.Service接口的模拟服务
type MockService struct{}

// 实现AirdropBnb方法
func (m *MockService) AirdropBnb(ctx context.Context, params service.AirdropParams) error {
	return nil
}

// 实现AirdropERC20方法
func (m *MockService) AirdropERC20(ctx context.Context, params service.AirdropParams) error {
	return nil
}

// 实现AirdropSetGov方法
func (m *MockService) AirdropSetGov(ctx context.Context, params service.AirdropSetGovParams) error {
	return nil
}

// 实现AirdropGov方法
func (m *MockService) AirdropGov(ctx context.Context) (string, error) {
	return "0x1234567890123456789012345678901234567890", nil
}

// 实现其他需要的方法
func (m *MockService) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*models.Block, error) {
	return &models.Block{}, nil
}

func (m *MockService) GetBlockByHash(ctx context.Context, blockHash string) (*models.Block, error) {
	return &models.Block{}, nil
}

func (m *MockService) SaveBlock(ctx context.Context, block *models.Block) error {
	return nil
}

func (m *MockService) GetLatestBlock(ctx context.Context) (*models.Block, error) {
	return &models.Block{}, nil
}

func (m *MockService) ERC20Allowance(ctx context.Context, params service.ERC20AllowanceParams) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *MockService) ERC20Approve(ctx context.Context, params service.ERC20ApproveParams) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *MockService) ERC20Transfer(ctx context.Context, params service.ERC20TransferParams) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *MockService) ERC20TransferFrom(ctx context.Context, params service.ERC20TransferFromParams) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *MockService) ERC20Balance(ctx context.Context, params service.ERC20BalanceParams) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *MockService) ERC20TotalSupply(ctx context.Context, params service.ERC20ContractParams) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *MockService) ERC20TokenInfo(ctx context.Context, params service.ERC20ContractParams) (*models.ERC20TokenInfo, error) {
	return &models.ERC20TokenInfo{}, nil
}

func main() {
	fmt.Println("===== 路由诊断工具 =====")

	// 1. 打印路由常量定义，验证它们是否正确
	fmt.Println("\n1. 验证路由常量定义:")
	fmt.Printf("   健康检查路径: %s\n", router.HealthPath)
	fmt.Printf("   设置空投合约地址: %s\n", router.AIRDROP_SET_GOV)
	fmt.Printf("   查询空投合约地址: %s\n", router.AIRDROP_GOV)
	fmt.Printf("   BNB空投: %s\n", router.AIRDROP_BNB)
	fmt.Printf("   ERC20空投: %s\n", router.AIRDROP_ERC20)

	// 2. 创建并测试一个与应用相同的路由实例
	fmt.Println("\n2. 创建并测试路由实例:")
	cfg := config.Config{}
	httpSrvCfg := config.HTTPServerConfig{}

	// 创建服务实例
	var svc service.Service = &MockService{}

	// 初始化路由
	r := router.InitRouter(httpSrvCfg, &cfg, svc)

	// 打印路由结构信息
	fmt.Printf("   路由类型: %v\n", reflect.TypeOf(r))

	// 3. 测试所有预期的路由
	fmt.Println("\n3. 测试所有预期路由:")

	testPaths := []struct {
		path   string
		method string
	}{{
		path:   router.HealthPath,
		method: "GET",
	}, {
		path:   router.AIRDROP_GOV,
		method: "GET",
	}, {
		path:   router.AIRDROP_SET_GOV,
		method: "POST",
	}, {
		path:   router.AIRDROP_BNB,
		method: "GET",
	}, {
		path:   router.AIRDROP_ERC20,
		method: "POST",
	}, {
		path:   "/api/nonexistent",
		method: "GET",
	}}

	// 测试每个路径
	for _, test := range testPaths {
		// 创建测试请求
		req := httptest.NewRequest(test.method, test.path, nil)
		rec := httptest.NewRecorder()

		// 调用路由
		r.ServeHTTP(rec, req)

		// 检查响应
		statusCode := rec.Code
		matched := statusCode != http.StatusNotFound
		statusText := "404 Not Found"
		if matched {
			statusText = "200 OK"
		}

		fmt.Printf("   %s %s: %s\n", test.method, test.path, statusText)

		// 如果匹配，打印处理器信息
		if matched && test.path != router.HealthPath {
			// 由于我们无法直接获取chi路由的处理器信息，我们只能确认路径是否匹配
			fmt.Printf("     ✓ 路径已匹配，但无法直接获取处理器信息\n")
		}
	}

	// 4. 创建一个简化版的路由进行对比测试
	fmt.Println("\n4. 简化版路由对比测试:")
	simpleRouter := chi.NewRouter()
	simpleRouter.Get(router.HealthPath, mockHandler)
	simpleRouter.Get(router.AIRDROP_GOV, mockHandler)
	simpleRouter.Post(router.AIRDROP_SET_GOV, mockHandler)
	simpleRouter.Get(router.AIRDROP_BNB, mockHandler)
	simpleRouter.Post(router.AIRDROP_ERC20, mockHandler)

	// 测试简化版路由
	simplereq := httptest.NewRequest("GET", router.AIRDROP_GOV, nil)
	simplerec := httptest.NewRecorder()
	simpleRouter.ServeHTTP(simplerec, simplereq)

	statusText := "404 Not Found"
	if simplerec.Code == http.StatusOK {
		statusText = "200 OK (匹配成功)"
	}
	fmt.Printf("   简化版路由 GET %s: %s\n", router.AIRDROP_GOV, statusText)

	fmt.Println("\n诊断完成！")
}
