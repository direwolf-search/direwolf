package repository

import (
	"context"
	"direwolf/internal/domain/model/host"
)

type HostRepository interface {
	CreateHost(ctx context.Context, h *host.Host) error
	UpdateHostByID(ctx context.Context, id int64, fields map[string]interface{}) error
	UpdateHostByURL(ctx context.Context, url string, fields map[string]interface{}) error
	GetHostByID(ctx context.Context, id int64) (*host.Host, error)
	GetHostByFields(ctx context.Context, fields map[string]interface{}) ([]*host.Host, error)
	GetAllHosts(ctx context.Context) ([]*host.Host, error)
	DeleteHost(ctx context.Context, id int64) error
}
