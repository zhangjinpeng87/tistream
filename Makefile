.PHONY: build

SRC = $(shell find . -type f -name '*.go' -not -path "./proto/*")

build:
	@echo "Building..."
	@GO111MODULE=on go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go
	@echo "Build complete"

proto_gen:
	@echo "Generating go proto files..."
	@protoc \
		--proto_path=proto \
		--go_out=proto/go/tistreampb \
		--go_opt=paths=source_relative \
		--go-grpc_out=proto/go/tistreampb \
		--go-grpc_opt=paths=source_relative \
		proto/*.proto
	@echo "Go proto files generated"

fmt:
	@gofmt -l -w $(SRC)
