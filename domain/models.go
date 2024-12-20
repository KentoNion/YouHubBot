package domain

import (
	"github.com/pkg/errors"
	"strings"
)

type UserID string

type Link string

var ValiableSourceLinks = map[string]bool{
	"youtube.com/watch": true,
	"youtu.be":          true,
	"vk.com/video":      false,
	"vk.com":            false,
	"vkvideo.ru/video":  false,
}

var ErrTgWrongLink = errors.New("Wrong link format, does not contain 't.me/'")
var ErrSourceWrongLink = errors.New("Wrong link format, does not contain any form of right format link")

type TgChan struct {
	Name           string
	Link           Link
	SourceChanLink Link
}

func VerifyChannelLink(link Link) error {
	contains := strings.Contains(string(link), "t.me/")
	if !contains {
		return ErrTgWrongLink
	}
	return nil
}

func VerifySourceLink(link Link) error {
	if ValiableSourceLinks[string(link)] {
		return nil
	}
	return errors.New("Wrong link format, does not contain link")
}
