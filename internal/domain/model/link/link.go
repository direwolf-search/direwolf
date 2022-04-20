package link

import (
	"errors"
)

var ErrInvalidLink = errors.New("error of invalid link")

type Link struct {
	// ID is a DB identifier of Link
	ID int64
	// From is a URL of Link's parent host
	From string
	// Body is a <href> attr value
	Body string
	// Snippet is an inner content of <a> tag
	Snippet string
	// IsV3 is an onion url version's boolean flag
	IsV3 bool
}

// NewLink is a Link's constructor
func NewLink(from, body, snippet string, isV3 bool) *Link {
	return &Link{
		From:    from,
		Body:    body,
		Snippet: snippet,
		IsV3:    isV3,
	}
}

// GetID returns link's ID.
// Link implements model.IDEntityGetter interface
func (l *Link) GetID() int64 {
	return l.ID
}

func (l *Link) Map() map[string]interface{} {
	return map[string]interface{}{
		"id":      l.ID,
		"from":    l.From,
		"body":    l.Body,
		"snippet": l.Snippet,
		"is_v3":   l.IsV3,
	}
}

func FromMap(m map[string]interface{}) (*Link, error) {
	var (
		l = &Link{}
	)

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

	err := l.Validate()
	if err != nil {
		return nil, err
	}

	if v, ok := m["id"]; ok {
		if int64Val, ok := v.(int64); ok {
			l.ID = int64Val
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

	return l, nil
}

func (l *Link) Validate() error {
	if l.From == "" || l.Body == "" {
		return ErrInvalidLink
	}

	return nil
}
