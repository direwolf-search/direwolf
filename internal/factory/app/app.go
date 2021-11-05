package app

import (
	"context"
)

type App interface {
	Do(ctx context.Context)
}
