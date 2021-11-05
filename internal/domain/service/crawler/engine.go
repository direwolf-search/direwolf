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
	Visit(ctx context.Context, f func(ctx context.Context, entity interface{}) error, url string)
	VisitAll(ctx context.Context, f func(ctx context.Context, entity interface{}) error, urls ...string)
	SaveLink(ctx context.Context, f func(ctx context.Context, entity interface{}) error, l *link.Link) error
	SaveHost(ctx context.Context, f func(ctx context.Context, entity interface{}) error, h *host.Host) error
	SetQueue()
	Init( /*c config.Config*/ )
	GetName() string
}
