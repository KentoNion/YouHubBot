package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type DB struct {
	DbUser string `yaml:"db_user" env-required:"true"`
	DbPass string `yaml:"db_password" env-required:"true"`
	DbHost string `yaml:"db_host" env-required:"true"`
	DbPort string `yaml:"db_port"`
	DbSsl  string `yaml:"db_sslmode" env-required:"true"`
}

type APIKeys struct {
	TeleApiKey string `yaml:"tele_api_key" env-required:"true"`
	YtApiKey   string `yaml:"yt_api_key" env-required:"true"`
}

type Log struct {
	FilePath string `yaml:"logger_file_path"`
}

type Config struct {
	Env     string  `yaml:"env"`
	DB      DB      `yaml:"postgres_db"`
	APIKeys APIKeys `yaml:"API_keys"`
	Log     Log     `yaml:"logger"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "../config.yaml"
	}

	//проверка существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("cannot read config file")
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
