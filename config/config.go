package config

import "time"

// Config 包含应用的配置信息
type Config struct {
	TokenExpiration time.Duration
	SecretKey       string
}

var AppConfig = Config{
	TokenExpiration: time.Hour * 24, // 默认设置为 24 小时
	SecretKey:       "your",
}
