GOHOSTOS := $(shell go env GOHOSTOS)
GOPATH := $(shell go env GOPATH)
VERSION := $(shell git describe --tags --always)

ifeq ($(GOHOSTOS), windows)
	MKDIR = if not exist bin mkdir bin
	BINARY = bin/soraka-backend.exe
	Git_Bash = $(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES := $(shell $(Git_Bash) -c "find internal -name *.proto")
#	API_PROTO_FILES := $(shell $(Git_Bash) -c "find api -name *.proto")
else
	MKDIR = mkdir -p bin
	BINARY = bin/soraka-backend
	INTERNAL_PROTO_FILES := $(shell find internal -name *.proto)
#	API_PROTO_FILES := $(shell find api -name *.proto)
endif

.PHONY: init
# 初始化依赖工具
init:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: wire
# wire依赖注入
wire:
	cd ./cmd && wire .

.PHONY: config
# 生成 internal proto 文件
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)
.PHONY: swag
swag:
	swag init -g ./cmd/main.go -o ./docs

.PHONY: api
# 生成 api proto 文件
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
	       --go_out=paths=source_relative:./api \
	       --go-http_out=paths=source_relative:./api \
	       --go-grpc_out=paths=source_relative:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)

.PHONY: generate
# 自动生成代码 + mod tidy
generate:
	go generate ./...
	go mod tidy

.PHONY: build
# 构建后端可执行文件
build:
	if not exist bin mkdir bin
	go build -ldflags "-X main.Version=$(VERSION)" -o bin/soraka-backend.exe ./cmd
.PHONY: all
# 一键生成 proto 和构建程序
all:
	make api
	make config
	make wire
	make generate
	make build

.PHONY: clean
# 清理构建产物
clean:
	rm -rf bin

.PHONY: help
# 显示帮助
help:
	@echo ''
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
