package main

import (
	"YoutHubBot/gates/postgres"
	"YoutHubBot/gates/telegram"
	"YoutHubBot/internal/config"
	"YoutHubBot/internal/logger"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //драйвер postgres
	"github.com/pkg/errors"
	goose "github.com/pressly/goose/v3"
	"go.uber.org/zap"
	telebot "gopkg.in/telebot.v3"
	"os"
	"time"
)

func main() {
	//читаем конфиг
	Cfg, err := config.MustLoad()
	if err != nil {
		panic(err)
	}

	//регестрируем логгер
	log, err := logger.InitLogger(Cfg)
	if err != nil {
		panic(err)
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:       Cfg.APIKeys.TeleApiKey,
		Synchronous: true,
		Verbose:     false,
		OnError: func(err error, msg telebot.Context) {
			if err != nil {
				log.Error("failed to send message", err)
			}
		},
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		panic(err)
	}

	dbhost := os.Getenv("DB_HOST") // получение значение DB_HOST из среды, значение среды todo: прописать значение среды в docker-compose
	if dbhost == "" {
		dbhost = Cfg.DB.DbHost
	}
	//подключение к дб
	connstr := fmt.Sprintf("user=%s password=%s dbname=youtube_hub_bot host=%s sslmode=%s", Cfg.DB.DbUser, Cfg.DB.DbPass, dbhost, Cfg.DB.DbSsl)
	conn, err := sqlx.Connect("postgres", connstr) //драйвер и имя бд
	if err != nil {
		zap.Error(errors.Wrap(err, "failed to connect to database"))
		panic(err)
	}
	db := postgres.NewDB(conn)
	//накатываем миграцию
	err = goose.Up(conn.DB, "./gates\\postgres\\migrations")
	if err != nil {
		panic(err)
	}

	opts := &telegram.Opts{
		Log: nil, //todo fix
	}
	client := telegram.NewClient(bot, opts)

	ctx := context.Background() //контекст
}
