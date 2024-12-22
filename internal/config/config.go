package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type DB struct {
	DbUser string `yaml:"db_user" env-required:"true"`
	DbPass string `yaml:"db_pass" env-required:"true"`
	DbHost string `yaml:"db_host" env-required:"true"`
	DbPort string `yaml:"db_port"`
	DbSsl  string `yaml:"db_ssl" env-required:"true"`
}

type APIKeys struct {
	TeleApiKey string `yaml:"tele_api_key" env-required:"true"`
	YtApiKey   string `yaml:"yt_api_key" env-required:"true"`
}

type Log struct {
	ConfigPath string `yaml:"config_path"`
}

type Config struct {
	Env     string `yaml:"env"`
	DB      DB
	APIKeys APIKeys
	Log     Log
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	//проверка существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("cannot read config file")
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot read config file", err))
	}
	return &cfg
}
