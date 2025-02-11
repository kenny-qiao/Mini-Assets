package config

import (
	"time"
)

// Config 结构体包含了应用程序的所有配置
type Config struct {
	ServerAddress string        // 服务器地址和端口
	DatabaseURL   string        // 数据库连接URL
	JWTSecret     string        // JWT密钥
	TokenDuration time.Duration // JWT令牌有效期
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	return &Config{
		ServerAddress: ":8080",                                                     // 服务器监听在8080端口
		DatabaseURL:   "asset:asset1234@tcp(assets-db:3306)/assets?parseTime=true", // MySQL数据库连接字符串
		JWTSecret:     "your_secret_key",                                           // JWT密钥，应该保持保密
		TokenDuration: 24 * time.Hour,                                              // JWT令牌有效期为24小时
	}
}
