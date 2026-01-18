package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string `yaml:"port"`
	Mode         string `yaml:"mode"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DBName       string `yaml:"db_name"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	// 对称密钥方式（如果使用HS256）
	SecretKey string `yaml:"secret_key"`

	// RSA密钥对方式（如果使用RS256）
	PublicKeyLocation  string `yaml:"public_key_location"`  // 公钥文件路径
	PrivateKeyLocation string `yaml:"private_key_location"` // 私钥文件路径

	Expiration int `yaml:"expiration"` // 过期时间（小时）
}

var AppConfig *Config

// LoadConfig 加载配置文件
// 如果 configPath 为空，会根据 APP_ENV 环境变量自动选择配置文件
// APP_ENV=dev -> config.dev.yaml
// APP_ENV=test -> config.test.yaml
// APP_ENV=prod -> config.prod.yaml
// 如果 APP_ENV 未设置，默认使用 config.yaml
func LoadConfig(configPath string) error {
	if configPath == "" {
		env := os.Getenv("APP_ENV")
		if env != "" {
			configPath = fmt.Sprintf("config.%s.yaml", env)
		} else {
			configPath = "config.yaml"
		}
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// 解析YAML配置
	AppConfig = &Config{}
	if err := yaml.Unmarshal(data, AppConfig); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// 设置默认值
	setDefaults()

	log.Printf("Config loaded from %s", configPath)
	return nil
}

// setDefaults 设置默认值
func setDefaults() {
	if AppConfig.Server.Port == "" {
		AppConfig.Server.Port = "8080"
	}
	if AppConfig.Server.Mode == "" {
		AppConfig.Server.Mode = "debug"
	}
	if AppConfig.Server.ReadTimeout == 0 {
		AppConfig.Server.ReadTimeout = 30
	}
	if AppConfig.Server.WriteTimeout == 0 {
		AppConfig.Server.WriteTimeout = 30
	}
	if AppConfig.Database.Host == "" {
		AppConfig.Database.Host = "localhost"
	}
	if AppConfig.Database.Port == "" {
		AppConfig.Database.Port = "3306"
	}
	if AppConfig.Database.MaxOpenConns == 0 {
		AppConfig.Database.MaxOpenConns = 100
	}
	if AppConfig.Database.MaxIdleConns == 0 {
		AppConfig.Database.MaxIdleConns = 10
	}
	if AppConfig.Redis.Host == "" {
		AppConfig.Redis.Host = "localhost"
	}
	if AppConfig.Redis.Port == "" {
		AppConfig.Redis.Port = "6379"
	}
	if AppConfig.JWT.SecretKey == "" {
		AppConfig.JWT.SecretKey = "your-secret-key-change-in-production"
	}
	if AppConfig.JWT.Expiration == 0 {
		AppConfig.JWT.Expiration = 24 // 默认24小时
	}
}
