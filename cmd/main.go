package main

import (
	"go.uber.org/zap"
	"gopkg.in/telebot.v3"
	"os"
)

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	bot, err := telebot.NewBot(telebot.Settings{
		Token:       os.Getenv("BOT_TELEGRAM_TOKEN"),
		Synchronous: true,
		Verbose:     false,
		OnError: func(err error, msg telebot.Context) {
			if err != nil {
				log.Error("failed to send message", zap.Error(err))
			}
		},
	})

}
