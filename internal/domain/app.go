package domain

import (
	"context"
)

type App interface {
	Do(ctx context.Context)
}
