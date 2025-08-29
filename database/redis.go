package database

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/urfave/cli/v2"
	"go-contracts/config"
	"go-contracts/util"
	"time"
)

// Redis 连接池包装器
type Redis struct {
	Pool *redis.Pool
}

// New 从配置创建Redis连接池
func NewRedis(c *cli.Context, cfg *config.RedisConfig) (*Redis, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	pool := &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: cfg.IdleTimeout,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			opts := []redis.DialOption{
				redis.DialConnectTimeout(2 * time.Second),
				redis.DialReadTimeout(2 * time.Second),
				redis.DialWriteTimeout(2 * time.Second),
			}

			if cfg.Password != "" {
				opts = append(opts, redis.DialPassword(cfg.Password))
			}

			return redis.Dial("tcp", addr, opts...)
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if c.Context != nil {
				// 4. 使用上下文进行健康检查
				select {
				case <-c.Context.Done():
					return c.Context.Err()
				default:
					_, err := conn.Do("PING")
					return err
				}
			}
			_, err := conn.Do("PING")
			return err
		},
	}

	// 验证连接
	conn := pool.Get()
	defer conn.Close()
	if _, err := conn.Do("PING"); err != nil {
		pool.Close()
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	util.Log.Info("Redis 链接成功", "address", addr)
	return &Redis{Pool: pool}, nil
}

// Close 关闭连接池
func (r *Redis) Close() error {
	if r.Pool != nil {
		return r.Pool.Close()
	}
	return nil
}
