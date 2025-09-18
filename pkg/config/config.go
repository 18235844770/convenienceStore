package config

import (
	"fmt"

	"github.com/spf13/viper"

	"convenienceStore/pkg/payment"
)

// AppConfig contains the consolidated application configuration schema.
type AppConfig struct {
	Server  ServerConfig   `mapstructure:"server"`
	Logging LoggingConfig  `mapstructure:"logging"`
	Payment payment.Config `mapstructure:"payment"`
}

// ServerConfig defines HTTP server runtime options.
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// Address renders the listen address string expected by gin.Engine.Run.
func (s ServerConfig) Address() string {
	if s.Host == "" {
		return fmt.Sprintf(":%d", s.Port)
	}
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// LoggingConfig captures structured logging settings.
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load reads configuration values from disk into AppConfig.
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
