package errors

//go:generate protoc --proto_path=. --proto_path=../../../../pkg --my-ext_out=errors=true,paths=source_relative:. --go_out=paths=source_relative:. ./article.proto

