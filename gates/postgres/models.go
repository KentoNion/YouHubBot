package postgres

import (
	"YoutHubBot/domain"
)

type Admin struct {
	userID domain.UserID `db:"admin"`
	role   string        `db:"admin_role"`
}

type Source struct {
	name           string      `db:"tg_chan_name""`
	link           domain.Link `db:"tg_chan_link"`
	sourceChanLink domain.Link `db:"tg_chan_source"`
}
