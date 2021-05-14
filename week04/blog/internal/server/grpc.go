package server

import (
	v1 "blog/api/blog/v1"
	"blog/internal/conf"
	"blog/internal/service"
	"my/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, blog *service.BlogService) *grpc.Server {
	var opts = []grpc.ServerOption{
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1.RegisterBlogServiceServer(srv, blog)
	return srv
}
