package host

type Host struct {
	ID         int64
	URL        string
	H1         string
	Title      string
	Links      []map[string]interface{}
	Body       string
	Hash       string
	Text       string
	Status     bool
	HTTPStatus string
	// Keywords
	//Ports   []*Port
	//Server    string
	//Proto     string
}

// GetID returns host's ID.
// Host implements model.IDEntityGetter interface
func (h *Host) GetID() int64 {
	return h.ID
}
