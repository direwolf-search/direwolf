package link

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

func (l *Link) LinkToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":      l.ID,
		"from":    l.From,
		"body":    l.Body,
		"snippet": l.Snippet,
		"is_v3":   l.IsV3,
	}
}
