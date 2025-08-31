package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go-contracts/config"
	"go-contracts/service"
	"net/http"
	"time"
)

const (
	// 健康检查路径
	HealthPath = "/api/health"

	// ERC20相关API路由
	ERC20_ALLOWANCE     = "/api/erc20/allowance"
	ERC20_APPROVE       = "/api/erc20/approve"
	ERC20_TRANSFER      = "/api/erc20/transfer"
	ERC20_TRANSFER_FROM = "/api/erc20/transfer_from"
	ERC20_BALANCE       = "/api/erc20/balance"
	ERC20_TOTAL_SUPPLY  = "/api/erc20/total_supply"
	ERC20_TOKEN_INFO    = "/api/erc20/token_info"

	// 空投相关路由
	AIRDROP_BNB   = "/api/airdrop_bnb"
	AIRDROP_ERC20 = "/api/airdrop_erc20"
)

func InitRouter(conf config.HTTPServerConfig, cfg *config.Config, svc service.Service) *chi.Mux {
	// 1. 创建验证器实例
	//	v := new(service.Validator)
	// 2. 创建业务服务实例
	//	svc := service.New(v, a.db.DepositTokens)
	// 3. 创建 chi 路由实例

	router := chi.NewRouter()
	// 创建路由处理器
	h := NewRoutes(router, svc)
	// 4. 注册中间件（与示例保持一致并添加新中间件）
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(10 * time.Second))
	router.Use(middleware.Heartbeat(HealthPath))
	// 5. 注册基础路由
	router.Get(HealthPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	
	
	// 注册空投相关路由
	router.Get(AIRDROP_BNB, h.AirdropBnb)// BNB空投
	router.Post(AIRDROP_ERC20, h.AirdropERC20)// ERC20空投

	// 注册ERC20相关路由
	router.Post(ERC20_ALLOWANCE, h.ERC20Allowance)// 查询授权
	router.Post(ERC20_APPROVE, h.ERC20Approve)// 授权
	router.Post(ERC20_TRANSFER, h.ERC20Transfer)// 转账
	router.Post(ERC20_TRANSFER_FROM, h.ERC20TransferFrom)// 从授权地址转账
	router.Post(ERC20_BALANCE, h.ERC20Balance)// 查询余额
	router.Post(ERC20_TOTAL_SUPPLY, h.ERC20TotalSupply)// 查询总供应量
	router.Post(ERC20_TOKEN_INFO, h.ERC20TokenInfo)// 查询代币信息
	

	return router
}
