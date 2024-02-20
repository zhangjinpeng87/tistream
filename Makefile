.PHONY: build

build:
	@echo "Building..."
	@GO111MODULE=on go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

proto_gen:
	@echo "Generating go proto files..."
	@protoc \
		--proto_path=proto \
		--go_out=proto/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=proto/go \
		--go-grpc_opt=paths=source_relative \
		proto/*.proto
	@echo "Go proto files generated"