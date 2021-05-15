.PHONY: generate build run test
NAME=myblog
VERSION=1.0.0

all: build 

URL:=https://raw.githubusercontent.com/googleapis/googleapis/master/google/api
P:=proto
PDIR:=pkg/third_party/google/api/
google-api-proto=\
	$(PDIR)/annotations.$(P) \
	$(PDIR)/http.$(P) \
	$(PDIR)/httpbody.$(P) \

$(google-api-proto):
	cd  ${PDIR} && ./down.bash "${URL}/$(notdir $@)"

down: $(google-api-proto)

protoc-gen=\
	protoc-gen-go_1.26 \
	protoc-gen-go-grpc_1.1 \
	protoc_3.14.0 \
	protoc-gen-my-ext_1.0.0 \

$(protoc-gen):
	$(eval tuple := $(subst _, ,$(@)))
	$(eval bin := $(word 1,$(tuple)))
	$(eval ver := $(word 2,$(tuple)))
	@which $(bin) 2>&1 > /dev/null || echo "need to install $@, please run 'make init'"
	@$(bin) --version|grep $(ver) 2>&1 >/dev/null && touch $(@) || echo "find $(bin) but diff with $@, please run 'make init'"

check_protoc-gen: $(protoc-gen)

wire:=wire_0.5.0
$(wire):
	@which wire 2>&1 >/dev/null && touch $(@) || echo "need to install wire. Please run 'make init'"

check_wire: $(wire)

check_exist:=\
	internal/data/ent/ent.go

$(check_exist):
	@[ -e $(@) ] || $(error "$(@) is missing, Please run 'make generate'")

check_exist: $(check_exist)

check: check_protoc-gen check_wire check_exist

init:
	go get github.com/google/wire/cmd/wire@v0.5.0
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	cd pkg/my/cmd/protoc-gen-my-ext/ && go install .

generate:
	cd pkg/my && go generate ./... && cd ../..
	go generate ./...
	wire ./...

build: check $(google-api-proto)
	@mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION) -X main.Name=$(NAME)" -o ./bin/ ./...

docker-up:
	docker-compose up

run:
	./bin/blog -conf ./configs/config.yaml

test:
	go test -v ./... -cover

clean:
	@rm -f $(protoc-gen) $(wire) bin/*

