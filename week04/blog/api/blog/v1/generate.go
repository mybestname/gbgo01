package v1

//go:generate protoc --proto_path=. --proto_path=../../../pkg --proto_path=../../../pkg/third_party --go-grpc_out=paths=source_relative:. --my-ext_out=errors=true,http=true,paths=source_relative:. --go_out=paths=source_relative:. ./blog.proto

