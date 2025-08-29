package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/urfave/cli/v2"
	"go-contracts/config"
	"go-contracts/controller/httputil"
	"go-contracts/database"
	"go-contracts/router"
	"go-contracts/service"
	"go-contracts/synchronizer/node"
	"go-contracts/util"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	HealthPath      = "/health"
	shutdownTimeout = 5 * time.Second
)

// ApiService HTTP 服务实现 cycle.Service 接口
type API struct {
	db        *database.DB    //  数据库连接
	redisPool *database.Redis // Redis
	apiServer *httputil.HTTPServer
	//	kafkaConsumer sarama.Consumer // Kafka 消费者（sarama）
	localCache *sync.Map // 本地缓存（sync.Map）
	stopped    atomic.Bool
	cfg        *config.Config
	router     *chi.Mux
}

// 创建 API 服务实例（业务入口）
func NewApi(c *cli.Context, cfg *config.Config) (*API, error) {
	api := &API{
		localCache: &sync.Map{},
	}
	if err := api.initFromConfig(c, cfg); err != nil {
		return nil, fmt.Errorf("Api创建失败: %w", err)
	}
	return api, nil
}
func (a *API) initFromConfig(c *cli.Context, cfg *config.Config) error {

	if err := a.initDB(c, cfg); err != nil {
		return fmt.Errorf("initDb初始化失败: %w", err)
	}
	if err := a.initRedis(c, cfg); err != nil {
		return fmt.Errorf("redis 初始化失败: %w", err)
	}
	// 创建请求参数验证器
	var v util.Validator
	
	// 初始化区块链客户端
	ethClient, err := node.DialEthClient(*c, config.RAW_URL)
	if err != nil {
		return fmt.Errorf("连接区块链节点失败: %w", err)
	}
	// 创建业务服务实例，传入区块对应链信息
	svc := service.New(v,  ethClient)
	// 初始化路由
	a.router = router.InitRouter(cfg.HTTPServer, cfg, svc)

	// 启动服务器
	if err := a.startServer(cfg.HTTPServer); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

// 初始化数据库（从命令行上下文）
func (a *API) initDB(c *cli.Context, cfg *config.Config) error {
	db, err := database.NewDb(c.Context, &cfg.MasterDB)
	if err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}
	a.db = db
	return nil
}

// 初始化redis
func (a *API) initRedis(c *cli.Context, cfg *config.Config) error {
	pool, err := database.NewRedis(c, &cfg.Redis)
	if err != nil {
		return err
	}
	a.redisPool = pool
	return nil
}

func (a *API) startServer(conf config.HTTPServerConfig) error {
	addr := net.JoinHostPort(conf.Host, fmt.Sprintf("%s", conf.Port))

	server, err := httputil.StartServerWithDefaults(addr, a.router)
	if err != nil {
		return fmt.Errorf("HTTP服务器启动失败: %w", err)
	}

	a.apiServer = server
	util.Log.Info("HTTP服务器已启动", "address", server.Addr().String())
	return nil
}

func (a *API) Start(ctx context.Context) error {
	if a.stopped.Load() {
		return errors.New("服务已停止，无法再次启动")
	}

	// 监听上下文取消信号
	go func() {
		<-ctx.Done()
		util.Log.Info("接收到停止信号，开始关闭服务...")
		if err := a.Stop(context.Background()); err != nil {
			util.Log.Error("服务关闭异常", "error", err)
		}
	}()

	util.Log.Info("API服务已启动")
	return nil
}

func (a *API) Stop(ctx context.Context) error {
	if a.stopped.Load() {
		return nil
	}

	var errs []error

	// 1. 关闭HTTP服务器
	if a.apiServer != nil {
		if err := a.apiServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("HTTP服务器关闭失败: %w", err))
		}
	}

	// 2. 关闭数据库连接
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			errs = append(errs, fmt.Errorf("数据库关闭失败: %w", err))
		}
	}

	// 3. 关闭Redis连接
	if a.redisPool != nil {
		if err := a.redisPool.Close(); err != nil {
			errs = append(errs, fmt.Errorf("Redis关闭失败: %w", err))
		}
	}

	// 4. 清理本地缓存
	if a.localCache != nil {
		a.localCache.Range(func(key, value interface{}) bool {
			a.localCache.Delete(key)
			return true
		})
	}

	a.stopped.Store(true)

	if len(errs) > 0 {
		return fmt.Errorf("服务关闭完成，但存在%d个错误: %w", len(errs), errors.Join(errs...))
	}

	util.Log.Info("API服务已正常关闭")
	return nil
}

func (a *API) Stopped() bool {
	return a.stopped.Load()
}
