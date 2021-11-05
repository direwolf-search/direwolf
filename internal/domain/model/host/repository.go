package host

import (
	"context"
)

type HostRepository interface {
	CreateHost(ctx context.Context, h *Host) error
	UpdateHostByID(ctx context.Context, id int64, fields map[string]interface{}) error
	UpdateHostByURL(ctx context.Context, url string, fields map[string]interface{}) error
	GetHostByID(ctx context.Context, id int64) (*Host, error)
	GetHostByFields(ctx context.Context, fields map[string]interface{}) ([]*Host, error)
	GetAllHosts(ctx context.Context) ([]*Host, error)
	DeleteHost(ctx context.Context, id int64) error
}
