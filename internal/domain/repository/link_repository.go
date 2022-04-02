package repository

import (
	"context"
	"direwolf/internal/domain/model/link"
)

type LinkRepository interface {
	CreateLink(context.Context, *link.Link) error
	UpdateLink(context.Context, int64, map[string]interface{}) error
	GetLinkByID(context.Context, int64) (*link.Link, error)
	GetLinkByFields(context.Context, map[string]interface{}) ([]*link.Link, error)
	GetLinksByHost(context.Context, int64) ([]*link.Link, error)
	GetAllLinks(context.Context) ([]*link.Link, error)
	DeleteLink(context.Context, int64) error
}
