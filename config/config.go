package config

import (
	"github.com/spf13/viper"
	"log"
	"strings"
)

var globalConfig Config

type Config struct {
	App struct {
		Name     string `mapstructure:"name"`
		Port     int    `mapstructure:"port"`
		LogLevel string `mapstructure:"log_level"`
	} `mapstructure:"app"`
	RabbitMQ struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		UserName string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"rmq"`
}

func LoadConfig() {
	config := viper.New()
	config.SetConfigName("app")
	config.SetConfigType("yaml")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
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
