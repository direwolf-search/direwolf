package models

import (
	"fmt"
	"github.com/uptrace/bun"
	"time"
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

func (l *Link) String() string {
	return fmt.Sprintf(" Link from: %s, body: %s", l.From, l.Body)
}
