package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Host struct {
	bun.BaseModel `bun:"hosts"`
	ID            int64     `bun:"id"`
	URL           string    `bun:"url"`
	H1            string    `bun:"h1,omitempty,nullzero"`
	Title         string    `bun:"title,omitempty,nullzero"`
	Links         []*Link   `bun:"-"`
	Hash          string    `bun:"hash,nullzero"`
	Text          string    `bun:"text,nullzero"`
	Status        bool      `bun:"status"`
	HTTPStatus    string    `bun:"http_status"`
	CreatedAt     time.Time `bun:"created_at,nullzero,hp:current_timestamp"`
	UpdatedAt     time.Time `bun:"updated_at,hp:current_timestamp"`
}

func NewHost() *Host {
	return &Host{}
}
