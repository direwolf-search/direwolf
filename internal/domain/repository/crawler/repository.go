package crawler

import (
	"context"
)

type Repository interface {
	GetAll(ctx context.Context) ([]string, error)
	//GetLinksCollections(ctx context.Context) ([]string, error) // not in MVP
	//GetByUpdateFreq(ctx context.Context, freq string) ([]string, error) // not in MVP
}
