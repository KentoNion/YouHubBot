package storage

import (
	"YoutHubBot/domain"
	"time"
)

type admin struct {
	userID domain.UserID `db:"admin"`
	role   string        `db:"admin_role"`
}

type source struct {
	name           string      `db:"tg_chan_name""`
	link           domain.Link `db:"tg_chan_link"`
	sourceChanLink domain.Link `db:"tg_chan_source"`
}

type articles struct {
	postID         int         `db:"post_ID"`
	tgChanName     string      `db:"tg_chan_name"`
	title          string      `db:"title"`
	postSourceLink domain.Link `db:"post_source_link"`
	publishedAt    time.Time   `db:"published_at"`
	createdAt      time.Time   `db:"created_at"`
	postedAt       time.Time   `db:"posted_at"`
}
