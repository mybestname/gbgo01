package http

import (
	"context"
	"my/transport"
	"net/http"
)

// ServerInfo represent HTTP server information.
type ServerInfo struct {
	Request  *http.Request
	Response http.ResponseWriter
}

type serverKey struct{}

// newServerContext returns a new http server Context that carries serverinfo.
func newServerContext(ctx context.Context, info ServerInfo) context.Context {
	ctx = transport.NewContext(ctx, transport.Transport{Kind: transport.HTTP})
	return context.WithValue(ctx, serverKey{}, info)
}