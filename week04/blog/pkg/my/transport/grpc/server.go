package grpc

import (
	"fmt"
	"net"
	"time"

	"my/api"
	"my/metadata"
	"my/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Server is a gRPC server wrapper.
type Server struct {
	*grpc.Server
	lis        net.Listener
	network    string
	address    string
	timeout    time.Duration
	log        *log.Helper
	grpcOpts   []grpc.ServerOption
	health     *health.Server
	metaServer *metadata.Server
}

// NewServer creates a gRPC server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		network: "tcp",
		address: ":0",
		timeout: time.Second,
		log:     log.NewHelper(loggerName, log.DefaultLogger),
		health: health.NewServer(),
	}
	for _, o := range opts {
		o(srv)
	}
	var grpcOpts []grpc.ServerOption
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	// grpc health register
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	// api metadata register
	api.RegisterMetadataServer(srv.Server, srv.metaServer)
	// reflection register
	reflection.Register(srv.Server)
	return srv
}


//
// Impl of the transport server interface
//

// Endpoint return a grpc address with port to registry endpoint.
// examples:
//   grpc://127.0.0.1:9000
func (s *Server) Endpoint() (string, error) {
	addr, port, err := net.SplitHostPort(s.address)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("grpc://%s", net.JoinHostPort(addr, port)), nil
}

// Start start the gRPC server.
func (s *Server) Start() error {
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		return err
	}
	s.lis = lis
	s.log.Infof("[gRPC] server listening on: %s", lis.Addr().String())
	s.health.Resume()
	return s.Serve(lis)
}

// Stop stop the gRPC server.
func (s *Server) Stop() error {
	s.GracefulStop()
	s.health.Shutdown()
	s.log.Info("[gRPC] server stopping")
	return nil
}