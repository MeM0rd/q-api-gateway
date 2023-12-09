server:
	go run cmd/q-api-gateway/main.go
mod:
	go mod tidy
pb:
	protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     pkg/pb/quote/quote.proto
