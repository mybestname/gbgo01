module backoff

go 1.16

require (
	google.golang.org/grpc v1.38.0
	grpcrand v0.0.0
)

replace grpcrand => ./grpcrand
