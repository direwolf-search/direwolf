package generic

import (
	"context"
)

type CommonRepository interface {
	Insert(ctx context.Context, entity interface{}) error
}
