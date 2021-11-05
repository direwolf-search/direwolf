package service

import (
	"context"
)

type Service interface {
	Do(ctx context.Context)
}
