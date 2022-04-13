package repository

import (
	"context"
)

type EngineRepository interface {
	Insert(ctx context.Context, entity map[string]interface{}) error
	Updated(ctx context.Context, url, md5hash string) (bool, error)
	Exists(ctx context.Context, url string) (bool, error)
	Update(ctx context.Context, entity map[string]interface{}) error
}
