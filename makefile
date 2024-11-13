all: generate build run

generate:
	@echo "Generating gRPC and grpc-gateway"
	protoc -I ./proto \
      --go_out ./pkg/api --go_opt paths=source_relative \
      --go-grpc_out ./pkg/api --go-grpc_opt paths=source_relative \
      --grpc-gateway_out ./pkg/api --grpc-gateway_opt paths=source_relative \
      ./proto/3_1/order.proto

build:
	@echo "building"
	go build cmd/main/main.go

run: build
	@echo "running"
	main.exe