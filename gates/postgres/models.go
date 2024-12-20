package postgres

import (
	"YoutHubBot/domain"
)

type TgChan struct {
	name           string
	link           domain.Link
	sourceChanLink domain.Link
}
