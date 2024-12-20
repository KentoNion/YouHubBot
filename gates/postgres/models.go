package postgres

import (
	"YoutHubBot/domain"
)

type TgChan struct {
	Name           string
	Link           domain.Link
	SourceChanLink domain.Link
}
