package storage

import (
	"YoutHubBot/domain"
	"database/sql"
	"time"
)

type admin struct {
	userID domain.UserID `db:"admin"`
	role   string        `db:"admin_role"`
}

type source struct {
	tg_name        string      `db:"tg_chan_name""`
	tg_link        domain.Link `db:"tg_chan_link"`
	sourceChanLink domain.Link `db:"tg_chan_source"`
}

type post struct {
	postID         int          `db:"ID"`
	SourceID       int          `db:"source_id"`
	tgChanName     string       `db:"tg_chan_name"`
	title          string       `db:"title"`
	postSourceLink domain.Link  `db:"post_source_link"`
	publishedAt    sql.NullTime `db:"published_at"`
	createdAt      time.Time    `db:"created_at"`
	postedAt       time.Time    `db:"posted_at"`
}
