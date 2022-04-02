package host

type Host struct {
	ID              int64
	URL             string
	Domain          string
	ContentType     string
	H1              string
	Title           string
	Links           []map[string]interface{}
	Meta            map[string]interface{}
	Body            string
	Hash            string
	Text            string // TODO:
	Status          bool
	HTTPStatus      string
	LinksCollection bool
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
