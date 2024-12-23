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

type user struct {
	domain.UserID
}

func NewClient(cli *tele.Bot, opts *Opts) *Client {
	svc := &Client{
		cli: cli,
		log: opts.Log,
	}
	cli.Handle(tele.OnText, svc.addReuploadChannel)
	cli.Handle("/addchanel", svc.addReuploadChannel)
	return svc
}

func (с *Client) hello(msg tele.Context) error {
	helloMessage := fmt.Sprintf("Wellcome to YouTube Hub Bot \n you may use this bot only if youre admin at this moment \n I have plans for the future to add limited permissions to download videos via dirrect messages to this bot in the future")
	msg.Reply(helloMessage)
	userId := msg.Sender().ID
	return nil
}

func (c *Client) addReuploadChannel(msg tele.Context) error {
	info := msg.Args()
	if len(info) != 3 {
		msg.Reply(fmt.Sprintf("Ввести нужно название канала и 2 ссылки, ссылку на YouTube/VK + TG чат/канал куда будет осуществляться презалив"))
	}
	channelName, channelLink, sourceLink := info[0], domain.Link(info[1]), domain.Link(info[2])
	err := domain.VerifyChannelLink(channelLink)
	if err != nil {
		msg.Reply(fmt.Sprintf("Неправильный формат ссылки, не содержит в себе ссылки на телеграм канал"))
		c.log.Info(fmt.Sprintf("failed to verify channel link: %s, User: %s", channelLink, msg.Sender().ID), zap.Error(err))
		return err
	}
	err = domain.VerifySourceLink(sourceLink)
	if err != nil {
		msg.Reply(fmt.Sprintf("Неправильный формат ссылки на источник, убедитесь что сообщение содержит в себе ссылку на видео"))
		c.log.Info(fmt.Sprintf("failed to verify source link: %s, User: %s", channelLink, msg.Sender().ID), zap.Error(err))
		return err
	}

	channel := domain.Source{
		Name:           channelName,
		Link:           channelLink,
		SourceChanLink: sourceLink,
	}
	//todo нужно как то проверить может ли бот писать в чат по ссылке
	cases.Subscrube(channel)
	if err != nil {
		msg.Reply(fmt.Sprint("Произошла ошибка"))
		c.log.Error(err)
		return err
	}
	return nil

}
