package link

import (
	"context"
)

type LinkRepository interface {
	CreateLink(context.Context, *Link) error
	UpdateLink(context.Context, int64, map[string]interface{}) error
	GetLinkByID(context.Context, int64) (*Link, error)
	GetLinkByFields(context.Context, map[string]interface{}) ([]*Link, error)
	GetLinksByHost(context.Context, int64) ([]*Link, error)
	GetAllLinks(context.Context) ([]*Link, error)
	DeleteLink(context.Context, int64) error
}
