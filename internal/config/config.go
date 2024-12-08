package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
	"os"
)

type Config struct {
	TeleApiKey string `yaml:"tele_api_key" env-required:"true"`
	YtApiKey   string `yaml:"yt_api_key" env-required:"true"`
	ConfigPath string `yaml:"config_path"`
}

func MustLoad() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	//проверка существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "config file not found at %s", configPath)
	}

	var cfg *Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read config file")
	}

	return cfg, nil
}
