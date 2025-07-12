.PHONY: help build run test clean docker-build docker-up docker-down docker-logs migrate

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
build: ## Build all Go binaries
	@echo "Building Go binaries..."
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker
	go build -o bin/rabbitmq-consumer ./cmd/rabbitmq_consumer
	go build -o bin/publisher ./cmd/publisher
	go build -o bin/subscriber ./cmd/subscriber

run: ## Run the application locally (requires services to be running)
	@echo "Running application..."
	go run cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker commands
docker-build: ## Build all Docker images
	@echo "Building Docker images..."
	docker compose build

docker-up: ## Start all services
	@echo "Starting services..."
	docker compose up -d

docker-down: ## Stop all services
	@echo "Stopping services..."
	docker compose down

docker-logs: ## Show logs for all services
	docker compose logs -f

docker-logs-api: ## Show API service logs
	docker compose logs -f api

docker-logs-worker: ## Show worker service logs
	docker compose logs -f workers

docker-logs-consumer: ## Show RabbitMQ consumer logs
	docker compose logs -f rabbitmq-consumer

# Database
migrate: ## Run database migrations
	@echo "Running database migrations..."
	docker compose run --rm migrate

db-reset: ## Reset database (WARNING: This will delete all data)
	@echo "WARNING: This will delete all data!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker compose down -v; \
		docker compose up -d postgres; \
		sleep 5; \
		make migrate; \
		echo "Database reset complete"; \
	else \
		echo "Database reset cancelled"; \
	fi

# Development utilities
dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	go mod download
	go mod tidy
	@echo "Development environment ready!"

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

format: ## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...

# Monitoring and debugging
status: ## Show service status
	@echo "Service Status:"
	docker compose ps

health: ## Check service health
	@echo "Checking service health..."
	@echo "API Health:"
	@curl -s http://localhost:8080/health || echo "API not responding"
	@echo "RabbitMQ Management:"
	@curl -s http://localhost:15672 || echo "RabbitMQ not responding"

# Testing utilities
test-mqtt: ## Test MQTT connection
	@echo "Testing MQTT connection..."
	mosquitto_pub -h localhost -t "test/topic" -m '{"test": "message"}'

test-geofence: ## Test geofence functionality
	@echo "Testing geofence entry..."
	mosquitto_pub -h localhost -t "vehicle/location" -m '{"vehicle_id":"TEST001","lat":-6.193125,"lng":106.820233}'

# Production
prod-build: ## Build production images
	@echo "Building production images..."
	docker compose -f docker-compose.yaml -f docker-compose.prod.yaml build

prod-deploy: ## Deploy to production
	@echo "Deploying to production..."
	docker compose -f docker-compose.yaml -f docker-compose.prod.yaml up -d

# Documentation
docs: ## Generate documentation
	@echo "Generating documentation..."
	@echo "Documentation is in README.md and docs/ directory"

# Performance testing
benchmark: ## Run performance benchmarks
	@echo "Running benchmarks..."
	go test -bench=. ./...

# Security
security-scan: ## Run security scan
	@echo "Running security scan..."
	gosec ./...

# Backup and restore
backup: ## Create database backup
	@echo "Creating database backup..."
	docker exec postgres pg_dump -U admin vehicle_tracker > backup_$(shell date +%Y%m%d_%H%M%S).sql

restore: ## Restore database from backup
	@echo "Restoring database from backup..."
	@read -p "Enter backup file name: " backup_file; \
	docker exec -i postgres psql -U admin vehicle_tracker < $$backup_file 