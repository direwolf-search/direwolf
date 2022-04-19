package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Host struct {
	bun.BaseModel  `bun:"hosts"`
	ID             int64                  `bun:"id"`
	URL            string                 `bun:"url"`
	Domain         string                 `bun:"domain"`
	ContentType    string                 `bun:"content_type"`
	H1             string                 `bun:"h1,omitempty,nullzero"`
	Title          string                 `bun:"title,omitempty,nullzero"`
	Links          []*Link                `bun:"-"`                       // TODO: ?????
	Meta           map[string]interface{} `bun:"meta,omitempty,nullzero"` // TODO: ???? flat?
	Md5hash        string                 `bun:"md5hash,nullzero"`
	Text           string                 `bun:"text,nullzero"`
	Status         bool                   `bun:"status"`
	HTTPStatus     string                 `bun:"http_status"`
	LinkCollection bool                   `bun:"link_collection"`
	CreatedAt      time.Time              `bun:"created_at,nullzero,hp:current_timestamp"` // TODO: current_timestamp???
	UpdatedAt      time.Time              `bun:"updated_at,hp:current_timestamp"`
	VisitedAt      time.Time              `bun:"updated_at,hp:current_timestamp"`
	VisitsNum      int                    `bun:"visits_num"`
}
