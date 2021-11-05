package repository

import (
	"context"

	"direwolf/internal/domain/model/host"
	"direwolf/internal/domain/model/link"
)

// Repository ...
type Repository interface {
	Insert(ctx context.Context, entity interface{}) error
	host.HostRepository
	link.LinkRepository
}
