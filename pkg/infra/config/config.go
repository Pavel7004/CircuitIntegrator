package config

import "github.com/spf13/viper"

type Config struct {
	Hostname string `mapstructure:"hostname"`
	Port     string `mapstructure:"port"`
}

var config *Config

func Read() error {
	viper.SetDefault("hostname", "localhost")
	viper.SetDefault("port", "8088")

	viper.AutomaticEnv()

	return viper.Unmarshal(config)
}

func init() {
	config = new(Config)
}

func Get() *Config {
	return config
}
