package repository

import (
	"context"

	"direwolf/internal/domain/model/host"
)

type CrawlerTaskPoolRepository interface {
	GetLinksCollections(ctx context.Context) ([]*host.Host, error)
	GetByUpdateFreq(ctx context.Context, freq string) ([]*host.Host, error)
}
