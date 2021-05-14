package transport

import "context"

// Server is transport server.
type Server interface {
	Endpoint() (string, error)
	Start() error
	Stop() error
}

// Kind defines the type of Transport
type Kind string

// Defines a set of transport kind
const (
	GRPC Kind = "gRPC"
	HTTP Kind = "HTTP"
)

// Transport is transport context value.
type Transport struct {
	Kind Kind
}

type transportKey struct{}

// NewContext returns a new Context that carries value.
func NewContext(ctx context.Context, tr Transport) context.Context {
	return context.WithValue(ctx, transportKey{}, tr)
}
