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
	HealthPath  = "/healthz"
	AIRDROP_BNB = "/api/airdrop_bnb"
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

	// 6. 注册业务路由（从 router 包引入）
	router.Get(AIRDROP_BNB, h.AirdropBnb)
	
	// 注册ERC20相关路由
	router.Post(ERC20_ALLOWANCE, h.ERC20Allowance)
	router.Post(ERC20_APPROVE, h.ERC20Approve)
	router.Post(ERC20_TRANSFER, h.ERC20Transfer)
	router.Post(ERC20_TRANSFER_FROM, h.ERC20TransferFrom)
	router.Post(ERC20_BALANCE, h.ERC20Balance)
	router.Post(ERC20_TOTAL_SUPPLY, h.ERC20TotalSupply)
	router.Post(ERC20_TOKEN_INFO, h.ERC20TokenInfo)

	return router
}
