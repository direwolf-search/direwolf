package crawler

import (
	"context"
	"net/http"

	"direwolf/internal/domain/model/host"
	"direwolf/internal/domain/model/link"
)

type Engine interface {
	GenerateRandomHeader() *http.Header
	SetRandomDelay()
	SetParallelism(workersNum int)
	SetTorGate(gate string)
	SetHTMLParser(parser interface{})
	Visit(ctx context.Context, url string)
	VisitAll(ctx context.Context, urls ...string)
	SaveLink(ctx context.Context, l *link.Link) error
	SaveHost(ctx context.Context, h *host.Host) error
	SetQueue()
	Init( /*c config.Config*/ )
	GetName() string
}
