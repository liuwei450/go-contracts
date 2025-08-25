package config

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 连接数据库
func ConnectDB(ctx context.Context, cfg *Config) (*gorm.DB, *sql.DB, error) {
	var (
		db  *gorm.DB
		err error
	)
	switch cfg.MasterDB.Driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.MasterDB.DSN()), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.MasterDB.DSN()), &gorm.Config{})
	default:
		return nil, nil, fmt.Errorf("不支持的数据库驱动: %s", cfg.MasterDB.Driver)
	}
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	return db, sqlDB, nil

}
