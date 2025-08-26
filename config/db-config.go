package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"path/filepath"
	"strings"
	"time"
)

type RedisConfig struct {
	Host        string        `yaml:"host"`
	Port        int           `yaml:"port"`
	Password    string        `yaml:"password"`
	MaxIdle     int           `yaml:"max_idle"`
	MaxActive   int           `yaml:"max_active"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}

type DBConfig struct {
	Driver          string        `yaml:"driver"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Name            string        `yaml:"name"`
	Config          string        `yaml:"config"`
	SSLMode         string        `yaml:"sslmode"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`    // 连接池最大空闲连接数
	MaxOpenConns    int           `yaml:"max_open_conns"`    // 连接池最大打开连接数
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"` // 连接最大存活时间（秒）
}

// HTTPServerConfig HTTP服务器配置
type HTTPServerConfig struct {
	Addr         string `yaml:"addr"`          // 监听地址（如 ":8080"）
	ReadTimeout  int    `yaml:"read_timeout"`  // 读取超时（秒）
	WriteTimeout int    `yaml:"write_timeout"` // 写入超时（秒）
	IdleTimeout  int    `yaml:"idle_timeout"`  // idle超时（秒）
}

type Config struct {
	MasterDB     DBConfig         `yaml:"masterdb"`     // 数据库配置
	MigrationDir string           `yaml:"migrationdir"` // 迁移文件目录
	HTTPServer   HTTPServerConfig `yaml:"httpserver"`   // HTTP服务器配置
	Redis        RedisConfig      `yaml:"redis"`        // Redis配置
	Kafka        KafkaConfig      `yaml:"kafka"`        // Kafka配置
}

const defaultConfigFileName = "config.yaml"

func LoadConfig(ctx *cli.Context) (*Config, error) {
	v := viper.New()
	var configPath string
	// 1. 配置文件路径
	cliConfigPath := ctx.String("config")
	//fmt.Println("数据库配置：" + cliConfigPath)
	if cliConfigPath != "" {
		configPath = cliConfigPath
		fmt.Printf("使用命令行指定的配置文件路径: %s\n", configPath)
	} else {
		configPath = defaultConfigFileName
		fmt.Printf("未指定 --config，使用默认配置文件路径: %s\n", configPath)
	}
	v.SetConfigFile(configPath)
	var configFileUsed string
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在：仅警告，继续使用默认值+环境变量
			fmt.Printf("未找到配置文件: %s，将使用默认值和环境变量\n", configPath)
		} else {
			// 文件存在但解析失败（如格式错误）：直接返回错误
			return nil, fmt.Errorf("配置文件读取失败（路径: %s）: %w", configPath, err)
		}
	} else {
		// 读取成功：记录实际使用的配置文件路径（绝对路径，方便调试）
		absPath, _ := filepath.Abs(v.ConfigFileUsed())
		configFileUsed = absPath
		fmt.Printf("成功读取配置文件: %s\n", configFileUsed)
	}
	// 2. 设置默认值
	v.SetDefault("masterdb.driver", "mysql")
	v.SetDefault("masterdb.port", 3306)
	v.SetDefault("masterdb.user", "root")
	v.SetDefault("masterdb.password", "123456")
	v.SetDefault("masterdb.sslmode", "")
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("kafka.brokers", []string{"localhost:9092"})
	v.SetDefault("migrationdir", "file://migrations")
	// ===== HTTP 服务器默认值 =====
	v.SetDefault("httpserver.addr", ":8080")     // 默认端口 8080
	v.SetDefault("httpserver.read_timeout", 10)  // 默认读取超时 10秒
	v.SetDefault("httpserver.write_timeout", 10) // 默认写入超时 10秒
	v.SetDefault("httpserver.idle_timeout", 30)  // 默认空闲超时 30秒

	// 3. 支持环境变量（自动大写并替换 . 为 _）
	v.SetEnvPrefix("APP") // 环境变量前缀 APP_
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 7. 映射配置到结构体（支持嵌套结构）
	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("配置映射到结构体失败: %w", err)
	}

	// 8. 验证关键配置（可选：确保必填项不为空）
	if cfg.MasterDB.Name == "" {
		return nil, fmt.Errorf("配置错误：masterdb.name（数据库名称）未设置，请检查配置文件或环境变量")
	}

	return cfg, nil
}

// DSN 生成数据库连接字符串
func (d DBConfig) DSN() string {
	switch d.Driver {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			d.User, d.Password, d.Host, d.Port, d.Name, d.Config)
	default:
		return fmt.Sprintf("不支持的数据库驱动: %s", d.Driver)
	}
}
