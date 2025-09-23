package config

import (
	"time"

	"github.com/spf13/viper"
)

type RateLimitConfig struct {
	Enabled bool          `mapstructure:"enabled"`
	Limit   int           `mapstructure:"limit"`
	Window  time.Duration `mapstructure:"window"`
}

type ServiceConfig struct {
	HTTP RateLimitConfig `mapstructure:"http"`
	GRPC RateLimitConfig `mapstructure:"grpc"`
}

type GlobalConfig struct {
	ServiceName string                   `mapstructure:"service_name"`
	RateLimits  map[string]ServiceConfig `mapstructure:"rate_limits"`
}

var AppConfig GlobalConfig

func LoadConfig(path string) error {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	return v.Unmarshal(&AppConfig)
}
