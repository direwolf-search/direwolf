package repository

import (
	"context"
)

type CrawlerRepository interface {
	GetAll(ctx context.Context) ([]string, error)
}
