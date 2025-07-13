.PHONY: help build run test clean docker-build docker-up docker-down docker-logs migrate

# Show this help message
help:
	@echo 'Usage: make [target]'
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build Go binaries
build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker
	go build -o bin/rabbitmq-consumer ./cmd/rabbitmq_consumer
	go build -o bin/publisher ./cmd/publisher
	go build -o bin/subscriber ./cmd/subscriber

# Run the application locally (requires services to be running)
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detector
test-race:
	go test -race ./...

# Run only delivery/http tests
test-http:
	go test -v ./internal/delivery/http/

# Run integration tests (requires services running)
test-integration:
	go test -v ./tests/integration/

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Docker commands
# Build all Docker images
docker-build:
	docker-compose build

# Start all services
docker-up:
	docker-compose up -d

# Stop all services
docker-down:
	docker-compose down

# Show logs for all services
docker-logs:
	docker-compose logs -f

# Show API service logs
docker-logs-api:
	docker-compose logs -f api

# Show worker service logs
docker-logs-worker:
	docker-compose logs -f workers

# Show RabbitMQ consumer logs
docker-logs-consumer:
	docker-compose logs -f rabbitmq-consumer

# Database
# Run database migrations
migrate:
	docker-compose run --rm migrate

# Reset database (WARNING: This will delete all data)
db-reset:
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker-compose down -v; \
		docker-compose up -d postgres; \
		sleep 5; \
		make migrate; \
		echo "Database reset complete"; \
	else \
		echo "Database reset cancelled"; \
	fi

# Development utilities
dev-setup:
	go mod download
	go mod tidy

# Monitoring and debugging
# Show service status
status:
	docker-compose ps

# Check service health
health:
	@echo "API Health:"
	@curl -s http://localhost:8080/healthz || echo "API not responding"
	@echo "RabbitMQ Management:"
	@curl -s http://localhost:15673 || echo "RabbitMQ not responding"

# Testing utilities
# Test MQTT connection
test-mqtt:
	mosquitto_pub -h localhost -t "test/topic" -m '{"test": "message"}'

# Test geofence functionality
test-geofence:
	mosquitto_pub -h localhost -t "vehicle/location" -m '{"vehicle_id":"TEST001","lat":-6.193125,"lng":106.820233}'

# Documentation
docs:
	@echo "Documentation is in README.md and docs/ directory"
