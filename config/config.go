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
		Host     string
		Port     int
		Name     string
		User     string
		Password string
	}

	ServerConfig struct {
		Name         string
		Version      string
		RPCPort      string
		RESTPort     string
		Debug        bool
		Environment  string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}

	AuthenticationConfig struct {
		Key       string
		SecretKey string
		SaltKey   string
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
