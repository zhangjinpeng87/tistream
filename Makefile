.PHONY: build

SRC = $(shell find . -type f -name '*.go' -not -path "./proto/*")

build:
	@echo "Building..."
	@GO111MODULE=on go build -o bin/meta-server cmd/meta-server/main.go;
	@GO111MODULE=on go build -o bin/api-server cmd/api-server/main.go;
	@GO111MODULE=on go build -o bin/dispatcher cmd/dispatcher/main.go;
	@GO111MODULE=on go build -o bin/remote-sorter cmd/remote-sorter/main.go;
	@GO111MODULE=on go build -o bin/schema-registry cmd/schema-registry/main.go;
	@echo "Build complete"

gen_proto:
	@echo "Generating go proto files..."
	@protoc \
		--proto_path=proto \
		--go_out=proto/go/tistreampb \
		--go_opt=paths=source_relative \
		--go-grpc_out=proto/go/tistreampb \
		--go-grpc_opt=paths=source_relative \
		proto/*.proto
	@echo "Go proto files generated"

clean_proto:
	@echo "Cleaning go proto files..."
	@rm proto/go/tistreampb/*
	@echo "Done"

fmt:
	@gofmt -l -w $(SRC)
