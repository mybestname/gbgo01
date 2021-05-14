module protoc-gen-my-ext

go 1.16

require (
	github.com/gorilla/mux v1.8.0
	google.golang.org/genproto v0.0.0-20210513213006-bf773b8c8384
	google.golang.org/grpc v1.37.1
	google.golang.org/protobuf v1.26.0
	my v1.0.0
)

replace my => ../../
