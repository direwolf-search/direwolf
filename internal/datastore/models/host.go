package models

import (
	"direwolf/internal/domain/model/host"
	"direwolf/internal/domain/model/link"
	"time"

	"github.com/uptrace/bun"
)

type Host struct {
	bun.BaseModel `bun:"hosts"`
	ID            int64                  `bun:"id"`
	URL           string                 `bun:"url"`
	Domain        string                 `bun:"domain"`
	ContentType   string                 `bun:"content_type"`
	H1            string                 `bun:"h1,omitempty,nullzero"`
	Title         string                 `bun:"title,omitempty,nullzero"`
	Links         []*Link                `bun:"links"`                   // TODO: ?????
	Meta          map[string]interface{} `bun:"meta,omitempty,nullzero"` // TODO: ???? flat?
	Md5hash       string                 `bun:"md5hash"`
	Text          string                 `bun:"text"`
	Status        bool                   `bun:"status"`
	HTTPStatus    string                 `bun:"http_status"`
	LinksNum      int                    `bun:"link_num"`
	CreatedAt     time.Time              `bun:"created_at,nullzero,hp:current_timestamp"` // TODO: current_timestamp???
	UpdatedAt     time.Time              `bun:"updated_at,hp:current_timestamp"`
	VisitedAt     time.Time              `bun:"updated_at,hp:current_timestamp"`
	VisitsNum     int                    `bun:"visits_num"`
}

func NewHostFromModel(h *host.Host) *Host {
	var links = make([]*Link, 0)

	for _, l := range h.Links {
		links = append(links, NewLinkFromModel(l))
	}

	return &Host{
		ID:          h.ID,
		URL:         h.URL,
		Domain:      h.Domain,
		ContentType: h.ContentType,
		H1:          h.H1,
		Title:       h.Title,
		Links:       links,
		Meta:        h.Meta,
		Md5hash:     h.MD5Hash,
		Text:        h.Text,
		Status:      h.Status,
		HTTPStatus:  h.HTTPStatus,
		LinksNum:    h.LinksNum,
	}
}

func (h *Host) ToModel() *host.Host {
	var links = make([]*link.Link, 0)

	for _, l := range h.Links {
		links = append(links, l.ToModel())
	}
	return &host.Host{
		ID:          h.ID,
		URL:         h.URL,
		Domain:      h.Domain,
		ContentType: h.ContentType,
		H1:          h.H1,
		Title:       h.Title,
		Links:       links,
		Meta:        h.Meta,
		MD5Hash:     h.Md5hash,
		Text:        h.Text,
		Status:      h.Status,
		HTTPStatus:  h.HTTPStatus,
		LinksNum:    h.LinksNum,
	}
}
