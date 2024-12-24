package domain

import (
	"github.com/pkg/errors"
	"strings"
	"time"
)

type UserID string

type Admin struct {
	UserID UserID
	Role   string
}

var ErrNotAdmin = errors.New("user not admin")

type Link string

var ValiableSourceLinks = map[string]bool{
	"youtube.com/watch": true,
	"youtu.be":          true,
	"vk.com/video":      false,
	"vk.com":            false,
	"vkvideo.ru/video":  false,
}

type Source struct {
	Name           string
	Link           Link
	SourceChanLink Link
}

var ErrTgWrongLink = errors.New("Wrong link format, does not contain 't.me/'")
var ErrSourceWrongLink = errors.New("Wrong link format, does not contain any form of right format link")
var ErrSourceAlreadyExist = errors.New("Source already exist")

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

type Post struct {
	postID         int
	tgChanName     string
	title          string
	postSourceLink Link
	publishedAt    time.Time
	createdAt      time.Time
	postedAt       time.Time
}
