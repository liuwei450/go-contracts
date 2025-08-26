package database

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // 适配器
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"strings"
)

// executeMigrations 执行迁移脚本
func ExecuteMigrations(ctx context.Context, dsn string, driver string, migrationDir string) error {
	// 根据 driver 决定适配器
	var dbURL string
	switch driver {
	case "postgres":
		dbURL = fmt.Sprintf("postgres://%s", strings.TrimPrefix(dsn, "postgres://"))
	case "mysql":
		dbURL = fmt.Sprintf("mysql://%s", strings.TrimPrefix(dsn, "mysql://"))
	default:
		return fmt.Errorf("不支持的数据库驱动: %s", driver)
	}

	m, err := migrate.New(migrationDir, dbURL)
	if err != nil {
		return fmt.Errorf("初始化迁移失败: %w", err)
	}

	// 执行迁移
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("迁移执行失败: %w", err)
	}
	return nil
}
