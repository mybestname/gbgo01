module balancer

go 1.16

require (
	github.com/go-kratos/kratos v1.0.0
	google.golang.org/grpc v1.29.1
	warden/metadata v0.0.0
)

replace warden/metadata => ./warden
