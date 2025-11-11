# Makefile для локальной разработки admin-gateway
# Используется для разработки одного сервиса

.PHONY: build test run tidy

# Update dependencies
tidy:
	@echo "📦 Updating dependencies..."
	@go mod tidy

# Build service
build:
	@echo "🔨 Building admin-gateway..."
	@go build -o bin/admin-gateway ./cmd/admin-gateway

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	@go test -cover ./...

# Run service locally (requires infrastructure to be running)
run:
	@echo "🚀 Running admin-gateway locally..."
	@go run ./cmd/admin-gateway

