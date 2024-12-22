package domain

import (
	"context"
)

type VideoProvider interface {
	GetVideo(ctx context.Context, link string) error
}
