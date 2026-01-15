package config

import (
	"log"
	"os"
	"strconv"

	"chat-server/pkg/types"
)

// LoadConfig 从环境变量加载配置
func LoadConfig() *types.ChatServerConfig {
	cfg := types.DefaultConfig()

	// 从环境变量读取配置
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		} else {
			log.Printf("无效的 PORT 值: %s, 使用默认值 %d", port, cfg.Port)
		}
	}

	if nodeEnv := os.Getenv("NODE_ENV"); nodeEnv == "production" {
		// 生产环境配置
		cfg.ReadTimeout = 30 * 60 // 30 minutes
		cfg.WriteTimeout = 30 * 60
	}

	return cfg
}
