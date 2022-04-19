package models

import (
	"fmt"
	"time"

	"github.com/uptrace/bun"

	"direwolf/internal/domain/model/link"
)

type Link struct {
	bun.BaseModel `bun:"links"`
	ID            int64     `bun:"id"`
	FromID        int64     `bun:"from_id"`
	From          string    `bun:"from"`
	Body          string    `bun:"body"`
	Snippet       string    `bun:"snippet"`
	IsV3          bool      `bun:"is_v3"`
	CreatedAt     time.Time `bun:"created_at"`
	UpdatedAt     time.Time `bun:"updated_at"`
}

func NewLinkFromModel(modelLink *link.Link) *Link {
	return &Link{
		ID:      modelLink.ID,
		From:    modelLink.From,
		Body:    modelLink.Body,
		Snippet: modelLink.Snippet,
		IsV3:    modelLink.IsV3,
	}
}

func NewLinkFromMap(m map[string]interface{}) *Link {
	var (
		l = &Link{}
	)

	if v, ok := m["id"]; ok {
		if int64Val, ok := v.(int64); ok {
			l.ID = int64Val
		}
	}
	if v, ok := m["from"]; ok {
		if stringVal, ok := v.(string); ok {
			l.From = stringVal
		}
	}
	if v, ok := m["body"]; ok {
		if stringVal, ok := v.(string); ok {
			l.Body = stringVal
		}
	}
	if v, ok := m["snippet"]; ok {
		if stringVal, ok := v.(string); ok {
			l.Snippet = stringVal
		}
	}
	if v, ok := m["is_v3"]; ok {
		if boolVal, ok := v.(bool); ok {
			l.IsV3 = boolVal
		}
	}

	return l
}

func (l *Link) ToModel() *link.Link {
	return &link.Link{
		ID:      l.ID,
		From:    l.From,
		Body:    l.Body,
		Snippet: l.Snippet,
		IsV3:    l.IsV3,
	}
}

func (l *Link) String() string {
	return fmt.Sprintf(" Link from: %s, body: %s", l.From, l.Body)
}
