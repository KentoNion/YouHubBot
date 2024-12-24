package main

import (
	"YoutHubBot/gates/storage"
	"YoutHubBot/gates/telegram"
	"YoutHubBot/internal/config"
	"YoutHubBot/internal/logger"
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
	Cfg := config.MustLoad()

	//регестрируем логгер
	log := logger.MustInitLogger(Cfg)
	log.Info("starting YoutHubBot")
	log.Debug("debug messages are enabled")

	//регестрируем бота
	bot, err := telebot.NewBot(telebot.Settings{
		Token:       Cfg.APIKeys.Telegram,
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

	// получение значение DB_HOST из среды, значение среды todo: прописать значение среды в docker-compose
	dbhost := os.Getenv("DB_HOST")
	if dbhost == "" {
		dbhost = Cfg.DB.Host
	}
	//подключение к дб
	connstr := fmt.Sprintf("user=%s password=%s dbname=youtube_hub_bot host=%s sslmode=%s", Cfg.DB.User, Cfg.DB.DbPass, dbhost, Cfg.DB.DbSsl)
	conn, err := sqlx.Connect("postgres", connstr) //драйвер и имя бд
	if err != nil {
		zap.Error(errors.Wrap(err, "failed to connect to database"))
		panic(err)
	}
	db := storage.NewDB(conn)

	//накатываем миграцию
	err = goose.Up(conn.DB, "./gates\\storage\\migrations")
	if err != nil {
		panic(err)
	}

	opts := &telegram.Opts{
		Log: nil, //todo fix
	}
	client := telegram.NewClient(bot, opts)

}
