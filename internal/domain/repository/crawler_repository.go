package repository

import (
	"context"
)

type CrawlerRepository interface {
	GetAll(ctx context.Context) ([]string, error)
	GetLinksCollections(ctx context.Context) ([]string, error)
	GetByUpdateFreq(ctx context.Context, freq string) ([]string, error)
}
