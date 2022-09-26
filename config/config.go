package config

import (
	"github.com/spf13/viper"
	"log"
)

var globalConfig Config

type Config struct {
	App struct {
		Name     string
		Port     int
		LogLevel string
	}
	RabbitMQ struct {
		Host     string
		Port     int
		UserName string
		Password string
	}
}

func LoadConfig() {
	config := viper.New()
	config.SetConfigName("app")
	config.SetConfigType("env")
	config.AutomaticEnv()

	config.AddConfigPath("./config/")
	config.AddConfigPath("../config/")

	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("read config error: %+v", err)
	}

	if err := config.Unmarshal(&globalConfig); err != nil {
		log.Fatalf("parse config error: %+v", err)
	}
}

func GetConfig() *Config {
	return &globalConfig
}
