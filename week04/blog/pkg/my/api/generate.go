package api
//go:generate protoc -I. -I../../third_party --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --my-ext_out=errors=true,http=true,paths=source_relative:. errors.proto metadata.proto


