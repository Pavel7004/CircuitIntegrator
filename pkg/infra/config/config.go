package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Hostname string        `mapstructure:"hostname"`
	Port     string        `mapstructure:"port"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

var config *Config

func Read() error {
	viper.SetDefault("hostname", "localhost")
	viper.SetDefault("port", "8000")
	viper.SetDefault("timeout", "10s")

	viper.AutomaticEnv()

	return viper.Unmarshal(config)
}

func init() {
	config = new(Config)
}

func Get() *Config {
	return config
}
