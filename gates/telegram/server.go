package telegram

import (
	"YoutHubBot/cases"
	"YoutHubBot/domain"
	"fmt"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
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
	cli.Handle(tele.OnText, svc.AddReuploadChannel())
	cli.Handle("/addchanel", svc.AddReuploadChannel)
	return svc
}

func (c *Client) AddReuploadChannel() func(msg tele.Context) {
	return func(msg tele.Context) error {
		info := msg.Args()
		if len(info) != 3 {
			msg.Reply(fmt.Sprintf("Ввести нужно название канала и 2 ссылки, ссылку на YouTube/VK + TG чат/канал куда будет осуществляться презалив"))
		}
		var channel domain.TgChan{
			name: info[0],
			link: info[1],
			sourceChanLink: info[2]
		}
		//todo нужно как то проверить может ли бот писать в чат по ссылке
		err := cases.Subscrube(channel)
		if err != nil {
			msg.Reply(fmt.Sprint("Произошла ошибка"))
			c.log.Error(err)
			return err
		}
		return nil
	}
}
