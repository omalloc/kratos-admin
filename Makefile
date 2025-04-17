GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always --abbrev=8 | sed 's/[^a-zA-Z0-9\.]/-/g')
GITHASH=$(shell git rev-parse HEAD)
APPNAME=$(shell go mod edit -print | head -n 1 | awk -F/ '{print $$3}')
Built:=$(shell date +%s)

ifeq ($(GOHOSTOS), windows)
  #the `find.exe` is different from `find` in bash/shell.
  #to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
  #changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
  #Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
  Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
  INTERNAL_PROTO_FILES=$(shell $(Git_Bash) -c "find internal -name *.proto")
  API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.proto")
  REMOTE_PROTO_FILES=$(shell $(Git_Bash) -c "go list -f '{{ .Dir }}' -m all | grep omalloc/contrib")
else
  INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
  API_PROTO_FILES=$(shell find api -name *.proto)
  REMOTE_PROTO_FILES=$(shell go list -f '{{ .Dir }}' -m all | grep omalloc/contrib)
endif

.PHONY: init
# init env
init:
	go install github.com/mitchellh/gox@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/google/wire/cmd/wire@latest

.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
		--proto_path=./third_party \
		--proto_path=$(REMOTE_PROTO_FILES) \
		--go_out=paths=source_relative:./internal \
		$(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
		--proto_path=./third_party \
		--proto_path=$(REMOTE_PROTO_FILES) \
		--go_out=paths=source_relative:./api \
		--go-errors_out=paths=source_relative:. \
		--go-http_out=paths=source_relative:./api \
		--go-grpc_out=paths=source_relative:./api \
		--openapi_out=fq_schema_naming=true,naming=proto,default_response=false:. \
		$(API_PROTO_FILES)

.PHONY: api-all
# generate openapi-v2, openapi-v3 api proto
api-all:
	mkdir -p ./api/docs && protoc --proto_path=./api \
		--proto_path=./third_party \
		--go_out=paths=source_relative:./api \
		--go-errors_out=paths=source_relative:. \
		--go-http_out=paths=source_relative:./api \
		--go-grpc_out=paths=source_relative:./api \
		--openapiv2_out=./api/docs \
		--openapiv2_opt=logtostderr=true \
		--openapiv2_opt=json_names_for_fields=false \
		--openapiv2_opt=output_format=json \
		--openapi_out=fq_schema_naming=true,naming=proto,default_response=false:. \
		$(API_PROTO_FILES)


.PHONY: run
run:
	@env go run -ldflags=" \
-X main.Version=$(VERSION) \
-X main.GitHash=$(GITHASH) \
-X main.Name=$(APPNAME) \
-X main.Built=$(Built) \
" \
		./cmd/... --conf ./configs/

.PHONY: build
# build
build:
	mkdir -p bin/ && gox -osarch="linux/amd64" -ldflags="-w -extldflags=-static \
		-X main.Version=$(VERSION) \
		-X main.GitHash=$(GITHASH) \
		-X main.Name=$(APPNAME) \
		-X main.Built=$(Built)" -output=bin/server ./...

.PHONY: zip
# zip bin file
zip:
	@upx bin/server --best --lzma

.PHONY: generate
# generate
generate:
	go mod tidy
	go get github.com/google/wire/cmd/wire@latest
	go generate ./...

.PHONY: wire
# wire generate wire file.
wire:
	@wire ./...

.PHONY: all
# generate all
all:
	make api;
	make config;
	make generate;

# show help
help:
	@echo ' '
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ' '
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
			} \
		} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
