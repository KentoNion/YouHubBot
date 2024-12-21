package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
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

func MustLoad() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}

	//проверка существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "config file not found at %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read config file")
	}

	return &cfg, nil
}
