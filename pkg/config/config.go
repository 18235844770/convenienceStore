package config

import (
	"fmt"

	"github.com/spf13/viper"

	"convenienceStore/pkg/payment"
)

// AppConfig 包含整个应用的配置结构。
type AppConfig struct {
	Server  ServerConfig   `mapstructure:"server"`
	Logging LoggingConfig  `mapstructure:"logging"`
	Payment payment.Config `mapstructure:"payment"`
}

// ServerConfig 定义 HTTP 服务器的运行时选项。
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// Address 返回 gin.Engine.Run 所需的监听地址字符串。
func (s ServerConfig) Address() string {
	if s.Host == "" {
		return fmt.Sprintf(":%d", s.Port)
	}
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// LoggingConfig 描述结构化日志的相关设置。
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load 从磁盘读取配置并填充 AppConfig。
func Load(path string) (*AppConfig, error) {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
