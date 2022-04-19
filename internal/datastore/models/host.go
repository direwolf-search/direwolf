package models

import (
	"time"

	"github.com/uptrace/bun"

	"direwolf/internal/domain/model/host"
	"direwolf/internal/domain/model/link"
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
	MD5Hash       string                 `bun:"md5hash"`
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
		MD5Hash:     h.MD5Hash,
		Text:        h.Text,
		Status:      h.Status,
		HTTPStatus:  h.HTTPStatus,
		LinksNum:    h.LinksNum,
	}
}

func NewHostFromMap(m map[string]interface{}) *Host {
	var (
		h = &Host{
			Links: make([]*Link, 0),
			Meta:  make(map[string]interface{}),
		}
	)

	if v, ok := m["id"]; ok {
		if int64Val, ok := v.(int64); ok {
			h.ID = int64Val
		}
	}
	if v, ok := m["url"]; ok {
		if stringVal, ok := v.(string); ok {
			h.URL = stringVal
		}
	}
	if v, ok := m["domain"]; ok {
		if stringVal, ok := v.(string); ok {
			h.Domain = stringVal
		}
	}
	if v, ok := m["content_type"]; ok {
		if stringVal, ok := v.(string); ok {
			h.ContentType = stringVal
		}
	}
	if v, ok := m["h1"]; ok {
		if stringVal, ok := v.(string); ok {
			h.H1 = stringVal
		}
	}
	if v, ok := m["title"]; ok {
		if stringVal, ok := v.(string); ok {
			h.Title = stringVal
		}
	}
	if v, ok := m["links"]; ok {
		if sl, ok := v.([]interface{}); ok {
			for _, interfaceVal := range sl {
				if mapVal, ok := interfaceVal.(map[string]interface{}); ok {
					l := NewLinkFromMap(mapVal)
					h.Links = append(h.Links, l)
				}
			}
		}
	}
	if v, ok := m["meta"]; ok {
		if mapVal, ok := v.(map[string]interface{}); ok {
			h.Meta = mapVal
		}
	}
	if v, ok := m["md5hash"]; ok {
		if stringVal, ok := v.(string); ok {
			h.MD5Hash = stringVal
		}
	}
	if v, ok := m["text"]; ok {
		if stringVal, ok := v.(string); ok {
			h.Text = stringVal
		}
	}
	if v, ok := m["status"]; ok {
		if boolVal, ok := v.(bool); ok {
			h.Status = boolVal
		}
	}
	if v, ok := m["http_status"]; ok {
		if stringVal, ok := v.(string); ok {
			h.HTTPStatus = stringVal
		}
	}
	if v, ok := m["links_num"]; ok {
		if intVal, ok := v.(int); ok {
			h.LinksNum = intVal
		}
	}

	return h
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
		MD5Hash:     h.MD5Hash,
		Text:        h.Text,
		Status:      h.Status,
		HTTPStatus:  h.HTTPStatus,
		LinksNum:    h.LinksNum,
	}
}
