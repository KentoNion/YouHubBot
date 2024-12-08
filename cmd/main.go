package main

import (
	"YoutHubBot/gates/postgres"
	"YoutHubBot/gates/telegram"
	"YoutHubBot/internal/config"
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
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
	conn, err := sqlx.Connect("postgres", "youtube_hub_bot.db") //драйвер и имя бд
	if err != nil {
		zap.Error(errors.Wrap(err, "failed to connect to database"))
		panic(err)
	}
	db := postgres.NewDB(conn) //подключение бд

	opts := &telegram.Opts{
		Log: log.With(zap.String("component", "telegram")),
	}
	client := telegram.NewClient(bot, opts)

	ctx := context.Background() //контекст
}
