package crawler

import (
	"context"
	"net/http"
)

type Engine interface {
	GenerateRandomHeader() *http.Header
	SetRandomDelay()
	SetParallelism(workersNum int)
	SetTorGate(gate string)
	SetHTMLParser(parser interface{})
	Visit(ctx context.Context, url string)
	VisitAll(ctx context.Context, urls ...string)
	SaveLink(ctx context.Context, l map[string]interface{}) error
	SaveHost(ctx context.Context, h map[string]interface{}) error
	SetQueue()
	Init( /*c config.Config*/ )
	GetName() string
}
