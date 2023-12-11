server:
	go run cmd/q-api-gateway/main.go
mod:
	go mod tidy
pb:
	protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative     pkg/pb/quote/quote.proto
make migrate up:
	migrate -path db/migrations -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' up
make migrate down:
	migrate -path db/migrations -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' down
