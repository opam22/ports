generate:
	@protoc --proto_path=api/ --go_out=internal/ports/grpc --go_opt=paths=source_relative \
		--go-grpc_out=internal/ports/grpc --go-grpc_opt=paths=source_relative  api/ports.proto

run-ports:
	@go run cmd/ports/main.go

run-importer:
	@go run cmd/importer/main.go

test:
	@go test ./... -v

docker-run:
	@docker-compose up

docker-build:
	@docker-compose build

check:
	@go fmt ./... && go vet ./...