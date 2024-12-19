package main

import (
	"YoutHubBot/gates/postgres"
	"YoutHubBot/gates/telegram"
	"YoutHubBot/internal/config"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //драйвер postgres
	"github.com/pkg/errors"
	goose "github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
	"os"
	"time"
)

func main() {
	Cfg, err := config.MustLoad() // ипортируем конфиг
	if err != nil {
		panic(err)
	}

	log, err := zap.NewDevelopment() //регестрируем логгер
	if err != nil {
		panic(err)
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:       Cfg.TeleApiKey,
		Synchronous: true,
		Verbose:     false,
		OnError: func(err error, msg telebot.Context) {
			if err != nil {
				log.Error("failed to send message", zap.Error(err))
			}
		},
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		panic(err)
	}

	dbhost := os.Getenv("DB_HOST") // получение значение DB_HOST из среды, значение среды todo: прописать значение среды в docker-compose
	if dbhost == "" {
		dbhost = Cfg.DbHost
	}
	//подключение к дб
	connstr := fmt.Sprintf("user=%s password=%s dbname=youtube_hub_bot host=%s sslmode=%s", Cfg.DbUser, Cfg.DbPass, dbhost, Cfg.DbSsl)
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
		Log: log.With(zap.String("component", "telegram")),
	}
	client := telegram.NewClient(bot, opts)

	ctx := context.Background() //контекст
}
