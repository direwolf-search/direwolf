// Package host defines a model of the host in Tor network
package host

import (
	"direwolf/internal/domain/model/link"
	"errors"
)

var ErrInvalidHost = errors.New("error of invalid host")

type Host struct {
	ID          int64
	URL         string
	Domain      string
	ContentType string
	H1          string
	Title       string
	Links       []*link.Link
	Meta        map[string]interface{}
	MD5Hash     string
	Text        string // TODO:
	Status      bool
	HTTPStatus  string
	LinksNum    int
	// Keywords
	//Ports   []*Port
	//Server    string
	//Proto     string
}

func (h *Host) Map() map[string]interface{} {
	var (
		ll = make([]map[string]interface{}, 0)
	)
	for _, l := range h.Links {
		linkMap := l.Map()
		ll = append(ll, linkMap)
	}
	return map[string]interface{}{
		"id":           h.ID,
		"url":          h.URL,
		"domain":       h.Domain,
		"content_type": h.ContentType,
		"h1":           h.H1,
		"title":        h.Title,
		"links":        ll,
		"meta":         h.Meta,
		"md5hash":      h.MD5Hash,
		"text":         h.Text,
		"status":       h.Status,
		"http_status":  h.HTTPStatus,
		"links_num":    h.LinksNum,
	}
}

func FromMap(m map[string]interface{}) (*Host, error) {
	var (
		h = &Host{
			Links: make([]*link.Link, 0),
			Meta:  make(map[string]interface{}),
		}
	)

	if v, ok := m["url"]; ok {
		if stringVal, ok := v.(string); ok {
			h.URL = stringVal
		}
	}

	err := h.Validate()
	if err != nil {
		return nil, err
	}

	if v, ok := m["id"]; ok {
		if int64Val, ok := v.(int64); ok {
			h.ID = int64Val
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
					l, _ := link.FromMap(mapVal)
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

	return h, nil
}

// GetID returns host's ID.
// Host implements model.IDEntityGetter interface
func (h *Host) GetID() int64 {
	return h.ID
}

func (h *Host) Validate() error {
	if h.URL == "" {
		return ErrInvalidHost
	}

	return nil
}
