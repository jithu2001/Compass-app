.PHONY: build run test clean docker-build docker-up docker-down

# Build the application
build:
	go build -o main ./cmd/api

# Run the application
run:
	go run ./cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f main
	go clean -cache

# Build Docker image
docker-build:
	docker-compose build

# Run with Docker
docker-up:
	docker-compose up -d

# Stop Docker containers
docker-down:
	docker-compose down

# View logs
logs:
	docker-compose logs -f

# Database migration (if needed in future)
migrate-up:
	@echo "Migrations are handled automatically by GORM"

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Download dependencies
deps:
	go mod download
	go mod tidy