module blog

go 1.16

require (
	entgo.io/ent v0.8.0
	github.com/go-redis/redis/extra/redisotel/v8 v8.8.2
	github.com/go-redis/redis/v8 v8.8.2
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/wire v0.5.0
	github.com/gorilla/mux v1.8.0
	go.opentelemetry.io/otel v0.19.0
	google.golang.org/genproto v0.0.0-20210513213006-bf773b8c8384
	google.golang.org/grpc v1.37.1
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v2 v2.4.0
	my v1.0.0
)

replace my => ./pkg/my
