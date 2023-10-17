package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Database       DatabaseConfig
		Server         ServerConfig
		Authentication AuthenticationConfig
	}

	DatabaseConfig struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	}

	ServerConfig struct {
		Name         string        `mapstructure:"name"`
		Version      string        `mapstructure:"version"`
		RPCPort      string        `mapstructure:"rpc_port"`
		RESTPort     string        `mapstructure:"rest_port"`
		Debug        bool          `mapstructure:"debug"`
		Environment  string        `mapstructure:"env"`
		ReadTimeout  time.Duration `mapstructure:"read_time_out"`
		WriteTimeout time.Duration `mapstructure:"write_time_out"`
	}

	AuthenticationConfig struct {
		Key       string `mapstructure:"key"`
		SecretKey string `mapstructure:"secret_key"`
		SaltKey   string `mapstructure:"salt_key"`
	}
)

func LoadConfig(env string) (Config, error) {
	viper.SetConfigFile(fmt.Sprintf("config/config-%s.yaml", env))

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return Config{}, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		return Config{}, err
	}

	return config, nil
}
