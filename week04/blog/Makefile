VERSION=1.0.0

.PHONY: generate build run test

URL=https://raw.githubusercontent.com/googleapis/googleapis/master/google/api
google-api:
	cd  pkg/third_party/google/api/ && ./down.bash \
		"${URL}/annotations.proto" "${URL}/http.proto" "${URL}/httpbody.proto"

init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	cd pkg/my/cmd/protoc-gen-my-ext/ && go install .

generate:
	go generate ./...

build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

run:
	cd cmd/blog/ && go run .

test:
	go test -v ./... -cover

