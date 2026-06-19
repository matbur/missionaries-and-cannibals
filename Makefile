MODULE := github.com/matbur/missionaries-and-cannibals
PROTOC_GEN_GO_VERSION := v1.36.11
PROTO := errors/errors.proto

.PHONY: help run build proto proto-install tidy clean

help:
	@echo "Targets:"
	@echo "  run           Run the solver"
	@echo "  build         Build the binary"
	@echo "  proto         Regenerate protobuf Go code"
	@echo "  proto-install Install protoc-gen-go ($(PROTOC_GEN_GO_VERSION))"
	@echo "  tidy          Run go mod tidy"
	@echo "  clean         Remove built binary"

run:
	go run main.go

build:
	go build -o missionaries-and-cannibals main.go

proto-install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)

proto: proto-install
	protoc -I=./errors --go_out=. --go_opt=module=$(MODULE) ./$(PROTO)

tidy:
	go mod tidy

clean:
	rm -f missionaries-and-cannibals
