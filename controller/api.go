package controller

import (
	"context"
	"errors"
	"fmt"
	"go-contracts/config"
	"go-contracts/database"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/urfave/cli/v2"
	"go-contracts/util"
)

const (
	defaultHTTPAddr         = ":8080"
	defaultHTTPReadTimeout  = 10 * time.Second
	defaultHTTPWriteTimeout = 10 * time.Second
	defaultHTTPIdleTimeout  = 30 * time.Second
	shutdownTimeout         = 5 * time.Second
)

// ApiService HTTP 服务实现 cycle.Service 接口
type API struct {
	server    *http.Server // HTTP 服务器实例
	db        *database.DB //  数据库连接
	redisPool *redis.Pool  // Redis 连接池（redigo）
	//	kafkaConsumer sarama.Consumer // Kafka 消费者（sarama）
	localCache *sync.Map // 本地缓存（sync.Map）
	stopped    atomic.Bool
	cfg        *config.Config
}

// 创建 API 服务实例（业务入口）
func NewApi(c *cli.Context, cfg *config.Config) (*API, error) {
	api := &API{
		localCache: &sync.Map{},
	}
	if err := api.initFromConfig(c); err != nil {
		return nil, fmt.Errorf("Api创建失败: %w", err)
	}
	return api, nil
}
func (a *API) initFromConfig(c *cli.Context) error {

	if err := a.initDB(c); err != nil {
		return fmt.Errorf("database initialization failed: %w", err)
	}
	if err := a.initRedis(); err != nil {
		return fmt.Errorf("redis initialization failed: %w", err)
	}

	if err := a.initHTTPServer(); err != nil {
		return fmt.Errorf("HTTP server initialization failed: %w", err)
	}

	return nil
}

// 初始化数据库（从命令行上下文）
func (a *API) initDB(c *cli.Context) error {
	db, err := database.FromCLIContext(context.Background(), c)
	if err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}
	a.db = db
	return nil
}

func (a *API) initRedis() error {
	redisAddr := fmt.Sprintf("%s:%d", a.cfg.Redis.Host, a.cfg.Redis.Port)
	a.redisPool = &redis.Pool{
		MaxIdle:     a.cfg.Redis.MaxIdle,
		MaxActive:   a.cfg.Redis.MaxActive,
		IdleTimeout: time.Duration(a.cfg.Redis.IdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			opts := []redis.DialOption{
				redis.DialConnectTimeout(2 * time.Second),
				redis.DialReadTimeout(2 * time.Second),
				redis.DialWriteTimeout(2 * time.Second),
			}
			if a.cfg.Redis.Password != "" {
				opts = append(opts, redis.DialPassword(a.cfg.Redis.Password))
			}
			return redis.Dial("tcp", redisAddr, opts...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// Verify connection
	conn := a.redisPool.Get()
	defer conn.Close()
	if _, err := conn.Do("PING"); err != nil {
		a.redisPool.Close()
		return fmt.Errorf("redis connection failed: %w", err)
	}

	util.Log.Info("Redis connection pool initialized")
	return nil
}

func (a *API) initHTTPServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	a.server = &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  defaultHTTPReadTimeout,
		WriteTimeout: defaultHTTPWriteTimeout,
		IdleTimeout:  defaultHTTPIdleTimeout,
	}

	util.Log.Info("HTTP server initialized")
	return nil
}

// Start 启动 HTTP 服务（阻塞方法，实现 cycle.Service 接口）
func (a *API) Start(ctx context.Context) error {
	if a.stopped.Load() {
		util.Log.Info("API 服务启动", "地址", a.server.Addr)
	}
	// 后台监听退出信号，关闭 HTTP 服务
	go func() {
		<-ctx.Done()
		util.Log.Info("开始关闭 HTTP 服务...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := a.server.Shutdown(shutdownCtx); err != nil {
			util.Log.Warn("HTTP 服务关闭异常", "error", err)
		}
	}()

	// 启动 HTTP 服务（阻塞）
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP 服务启动失败: %w", err)
	}
	util.Log.Info("HTTP 服务已正常关闭")
	return nil
}

// Stop 停止服务并清理所有资源（实现 cycle.Service 接口）
func (a *API) Stop(ctx context.Context) error {
	if a.stopped.Load() {
		return nil
	}
	var allErrors []error // 收集所有清理步骤的错误
	// ===== 1. 关闭 HTTP 服务（冗余检查，确保已关闭）=====
	if a.server != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		if err := a.server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			allErrors = append(allErrors, fmt.Errorf("HTTP 服务关闭失败: %w", err))
		}
		a.server = nil // 释放引用
	}

	// ===== 2. 关闭数据库连接 =====
	if err := a.db.Close(); err != nil {
		allErrors = append(allErrors, fmt.Errorf("关闭数据库失败: %w", err))
	}

	// ===== 3. 关闭 Redis 连接池 =====
	if a.redisPool != nil {
		util.Log.Info("关闭 Redis 连接池...")
		a.redisPool.Close() // 关闭所有连接
		a.redisPool = nil   // 释放引用
		util.Log.Info("Redis 连接池已关闭")
	}

	// ===== 4. 清理本地缓存 =====
	if a.localCache != nil {
		util.Log.Info("清理本地缓存...")
		a.localCache.Range(func(key, value interface{}) bool {
			a.localCache.Delete(key) // 清空所有键值对
			return true
		})
		a.localCache = nil // 释放引用
		util.Log.Info("本地缓存已清理")
	}

	// ===== 返回汇总错误 =====
	if len(allErrors) > (0) {
		return fmt.Errorf("资源清理完成，但存在 %d个错误: %w", len(allErrors), errors.Join(allErrors...))
	}

	util.Log.Info("API 服务所有资源已成功清理")
	return nil
}

func (a *API) Stopped() bool {
	return a.stopped.Load()
}
