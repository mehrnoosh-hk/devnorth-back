.PHONY: help build run test clean sqlc migrate-up migrate-down migrate-create

help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make sqlc           - Generate Go code from SQL queries"
	@echo "  make migrate-up     - Run database migrations up"
	@echo "  make migrate-down   - Rollback last migration"
	@echo "  make migrate-create - Create a new migration (use name=your_migration_name)"

build:
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go

run:
	@echo "Running application..."
	go run cmd/api/main.go

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/

sqlc:
	@echo "Generating SQLC code..."
	sqlc generate

migrate-up:
	@echo "Running migrations..."
	migrate -path db/migrations -database "$(DB_URL)" up

migrate-down:
	@echo "Rolling back migration..."
	migrate -path db/migrations -database "$(DB_URL)" down 1

migrate-create:
	@echo "Creating migration: $(name)"
	migrate create -ext sql -dir db/migrations -seq $(name)
