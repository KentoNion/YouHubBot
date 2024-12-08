package telegram

import (
	"YoutHubBot/domain"
	"fmt"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"strings"
)

type Client struct {
	cli *tele.Bot
	log *zap.Logger
}

type Opts struct {
	Provider domain.VideoProvider
	Log      *zap.Logger
}

func NewClient(cli *tele.Bot, opts *Opts) *Client {
	svc := &Client{
		cli: cli,
		log: opts.Log,
	}
	cli.Handle(tele.OnText, svc.addReuploadChannel())
	cli.Handle("/addchanel", svc.addReuploadChannel())
	return svc
}

func (c *Client) addReuploadChannel(msg string) func(msg tele.Context) error {
	return func(msg tele.Context) error {
		links := strings.Split(msg)
		if len(links) != 2 {
			msg.Reply(fmt.Sprintf("Ввести нужно 2 ссылки, ссылку на YouTube/VK + TG чат/канал куда будет осуществляться презалив"))
			return nil
		}
		sourceLink, tgChanLink := links[0], links[1]
		//todo нужно как то проверить может ли бот писать в чат по ссылке
	}
}
