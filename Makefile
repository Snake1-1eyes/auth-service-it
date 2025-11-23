include .env

APP_NAME := auth
BUILD_DIR := build
LOCAL_BIN := $(CURDIR)/bin
OUT_PATH := $(CURDIR)/pkg
GOOSE := $(LOCAL_BIN)/goose

update:
	go mod tidy

linter:
	golangci-lint run ./...

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/auth-api

start-grpc:
	$(BUILD_DIR)/$(APP_NAME)

run: update linter build start-grpc

kill-all:
	pkill -f "$(BUILD_DIR)/$(APP_NAME)$$"

bin-deps: export GOBIN := $(LOCAL_BIN)
bin-deps: export PROTOC_VERSION := protoc-31.1-linux-x86_32
bin-deps:
	mkdir -p $(LOCAL_BIN)
	curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v31.1/$(PROTOC_VERSION).zip
	unzip -o $(PROTOC_VERSION).zip -d $(LOCAL_BIN)
	rm $(PROTOC_VERSION).zip

	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

generate: export GOBIN := $(LOCAL_BIN)
generate:
	mkdir -p $(OUT_PATH)
	$(LOCAL_BIN)/bin/protoc --proto_path=api --proto_path=vendor.protogen \
		--go_out=$(OUT_PATH) --go_opt=paths=source_relative --plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
		--go-grpc_out=$(OUT_PATH) --go-grpc_opt=paths=source_relative --plugin protoc-gen-go-grpc="${GOBIN}/protoc-gen-go-grpc" \
		--validate_out="lang=go,paths=source_relative:$(OUT_PATH)" --plugin protoc-gen-validate=$(LOCAL_BIN)/protoc-gen-validate \
		--grpc-gateway_out=$(OUT_PATH) --grpc-gateway_opt=paths=source_relative --plugin protoc-gen-grpc-gateway=$(LOCAL_BIN)/protoc-gen-grpc-gateway \
		--openapiv2_out=$(OUT_PATH) --plugin=protoc-gen-openapiv2=$(LOCAL_BIN)/protoc-gen-openapiv2 \
		api/auth/auth.proto
	go mod tidy

.vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
	https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp && \
	cd vendor.protogen/tmp && \
	git sparse-checkout set --no-cone validate && \
	git checkout
	mkdir -p vendor.protogen/validate
	mv vendor.protogen/tmp/validate vendor.protogen/
	rm -rf vendor.protogen/tmp

.vendor-proto/google/api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
 		cd vendor.protogen/googleapis && \
		git sparse-checkout set --no-cone google/api && \
		git checkout
		mkdir -p vendor.protogen/google
		mv vendor.protogen/googleapis/google/api vendor.protogen/google
		rm -rf vendor.protogen/googleapis

.vendor-proto/protoc-gen-openapiv2/options:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
 		https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem && \
 		cd vendor.protogen/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
		mkdir -p vendor.protogen/protoc-gen-openapiv2
		mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options vendor.protogen/protoc-gen-openapiv2
		rm -rf vendor.protogen/grpc-ecosystem

up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

install-goose:
	mkdir -p $(LOCAL_BIN)
	@GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

goose-create: install-goose
	$(GOOSE) -dir ./migrations postgres "$(POSTGRES_DSN)" create $(NAME) sql

goose-up: install-goose
	$(GOOSE) -dir ./migrations postgres "$(POSTGRES_DSN)" up

goose-down: install-goose
	$(GOOSE) -dir ./migrations postgres "$(POSTGRES_DSN)" down

goose-status: install-goose
	$(GOOSE) -dir ./migrations postgres "$(POSTGRES_DSN)" status
