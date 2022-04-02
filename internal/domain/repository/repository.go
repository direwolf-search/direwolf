package repository

import (
	"context"
)

// Repository ...
type Repository interface {
	Insert(ctx context.Context, entity interface{}) error
	Updated(ctx context.Context, url, md5hash string) (bool, error)
	Exists(ctx context.Context, url string) (bool, error)
	Update(ctx context.Context, url string) error
	HostRepository
	LinkRepository
	CrawlerTaskPoolRepository
}
