package database

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"go-contracts/config"
	"go-contracts/util"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// DB 数据库连接实例（包装gorm.DB，对外提供统一接口）
type DB struct {
	*gorm.DB
}

// New 从config.DBConfig创建数据库连接
// 直接使用config包的DBConfig，避免重复定义配置结构
func New(ctx context.Context, cfg *config.DBConfig) (*DB, error) {
	// 1. 根据数据库驱动类型初始化GORM
	var dialector gorm.Dialector
	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.DSN()) // 使用config.DBConfig自带的DSN()方法
	case "postgres":
		dialector = postgres.Open(cfg.DSN())
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.Driver)
	}

	// 2. 初始化GORM连接
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 日志级别可配置
	})
	if err != nil {
		return nil, fmt.Errorf("GORM连接失败: %w", err)
	}

	// 3. 配置连接池（使用config.DBConfig的参数）
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层SQL连接失败: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime * time.Second)

	// 4. 验证连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库Ping失败: %w", err)
	}

	util.Log.Info("数据库连接成功", "driver", cfg.Driver, "host", cfg.Host, "dbname", cfg.Name)
	return &DB{db}, nil
}

// FromCLIContext 从命令行上下文创建数据库连接
// 流程：加载全局配置 → 提取MasterDB配置 → 初始化数据库
func FromCLIContext(ctx context.Context, c *cli.Context) (*DB, error) {
	// 1. 调用config包加载全局配置（已实现）
	globalCfg, err := config.LoadConfig(c)
	if err != nil {
		return nil, fmt.Errorf("加载全局配置失败: %w", err)
	}

	// 2. 直接使用globalCfg.MasterDB（config.DBConfig类型）
	return New(ctx, &globalCfg.MasterDB)
}

// Close 关闭数据库连接
func (d *DB) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("获取底层SQL连接失败: %w", err)
	}
	d.DB = nil
	return sqlDB.Close()
}
