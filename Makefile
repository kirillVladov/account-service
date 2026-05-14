MIGRATIONS_PATH=build/app/migrations

ifneq (,$(wildcard build/local/.env))
  include build/local/.env
  export
endif

run: 
	bash -c 'set -a; . ./build/local/.env; set +a; go run cmd/account-service/main.go run'
	
lint:
	@golangci-lint run

test:
	@go test ./... -v

up:
	docker-compose up

down:
	docker-compose down

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(POSTGRES_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(POSTGRES_URL)" down

GOPATH_BIN := $(shell go env GOPATH)/bin
PROTOC_GEN_GO := $(GOPATH_BIN)/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOPATH_BIN)/protoc-gen-go-grpc

.PHONY: install-protoc-gen
install-protoc-gen:
	@echo "Installing protoc plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: generate
generate: install-protoc-gen
	@echo "Generating gRPC code..."
	@protoc \
		--plugin=protoc-gen-go=$(PROTOC_GEN_GO) \
		--plugin=protoc-gen-go-grpc=$(PROTOC_GEN_GO_GRPC) \
		--go_out=. \
		--go-grpc_out=. \
		internal/docs/proto/*.proto
	@echo "✅ Done!"

.PHONY: clean
clean:
	@rm -f generated/*.pb.go
	@echo "Cleaned generated files"
