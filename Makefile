.PHONY: all build run dev clean test frontend backend docker

# Default target
all: build

# Build everything
build: frontend backend

# Build frontend
frontend:
	cd web && npm ci && npm run build

# Build backend
backend:
	go build -o bin/listenbucket ./cmd/listenbucket

# Run in development mode (requires frontend built)
run: backend
	./bin/listenbucket

# Development mode with hot reload
dev:
	@echo "Starting backend on :8080..."
	@echo "Starting frontend dev server on :5173..."
	@echo "Frontend will proxy API requests to backend"
	@trap 'kill %1 %2 2>/dev/null' EXIT; \
	go run ./cmd/listenbucket & \
	cd web && npm run dev & \
	wait

# Run frontend dev server only
dev-frontend:
	cd web && npm run dev

# Run backend only
dev-backend:
	go run ./cmd/listenbucket

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf web/dist/
	rm -rf web/node_modules/

# Run tests
test:
	go test -v ./...

# Build Docker image
docker:
	docker build -t listenbucket:latest .

# Run with Docker Compose
docker-up:
	docker-compose up -d

# Stop Docker Compose
docker-down:
	docker-compose down

# View Docker logs
docker-logs:
	docker-compose logs -f

# Install dependencies
deps:
	go mod download
	cd web && npm ci

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	go vet ./...
