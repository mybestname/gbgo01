package main

//go:generate go install .
//go:generate protoc --proto_path=. --proto_path=../../.. --proto_path=../../../third_party --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. --my-ext_out=errors=true,http=true,paths=source_relative:.  ./test/error_test.proto ./test/http_test.proto

