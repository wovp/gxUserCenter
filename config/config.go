package config

import "time"

// Config 包含应用的配置信息
type Config struct {
	// token的过期时间
	TokenExpiration time.Duration
	// token密钥
	SecretKey string
}

var AppConfig = Config{
	TokenExpiration: time.Hour * 24, // 默认设置为 24 小时
	SecretKey:       "your",
}
