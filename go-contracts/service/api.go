package service

import (
	"context"
	"errors"
	"fmt"
	"go-contracts/config"
	"gorm.io/driver/mysql"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"

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
	db        *gorm.DB     // GORM v2 数据库连接
	redisPool *redis.Pool  // Redis 连接池（redigo）
	//	kafkaConsumer sarama.Consumer // Kafka 消费者（sarama）
	localCache *sync.Map // 本地缓存（sync.Map）
	stopped    atomic.Bool
	cfg        *config.Config
}

// 创建 API 服务实例（业务入口）
func NewApiService(c *cli.Context) (*API, error) {
	api := &API{
		localCache: &sync.Map{},
	}
	if err := api.initFromConfig(c); err != nil {
		return nil, fmt.Errorf("Api创建失败: %w", err)
	}
	return api, nil
}
func (a *API) initFromConfig(c *cli.Context) error {
	cfg, err := a.loadConfig(c)
	if err != nil {
		return fmt.Errorf("数据库配置信息加载失败: %w", err)
	}
	a.cfg = cfg
	if err := a.initDB(); err != nil {
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

// 加载配置
func (a *API) loadConfig(c *cli.Context) (*config.Config, error) {
	cfg, err := config.LoadConfig(c)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}

// 数据库初始化
func (a *API) initDB() error {
	dsn := a.cfg.MasterDB.DSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying DB connection: %w", err)
	}
	// Configure connection pool
	sqlDB.SetMaxIdleConns(a.cfg.MasterDB.MaxIdleConns)
	sqlDB.SetMaxOpenConns(a.cfg.MasterDB.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(a.cfg.MasterDB.ConnMaxLifetime) * time.Second)

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	a.db = db
	util.Log.Info("Database connection established")
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
		Addr:         a.cfg.HTTPServer.Addr,
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
	if a.db != nil {
		util.Log.Info("关闭数据库连接...")
		sqlDB, err := a.db.DB()
		if err != nil {
			allErrors = append(allErrors, fmt.Errorf("获取数据库底层连接失败: %w", err))
		} else {
			// 设置关闭超时（等待活跃连接完成）
			sqlDB.SetConnMaxLifetime(0)
			sqlDB.SetMaxIdleConns(0)
			if err := sqlDB.Close(); err != nil {
				allErrors = append(allErrors, fmt.Errorf("数据库关闭失败: %w", err))
			} else {
				util.Log.Info("数据库连接已关闭")
			}
		}
		a.db = nil // 释放引用
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
