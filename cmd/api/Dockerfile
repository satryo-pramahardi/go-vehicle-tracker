# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the API binary
RUN go build -o api ./cmd/api

# Runtime stage
FROM alpine
WORKDIR /app

# Install ca-certificates for DB/Redis TLS if needed
RUN apk add --no-cache ca-certificates

# Copy binary
COPY --from=builder /app/api .

# Expose port
EXPOSE 8080

# Default command
CMD ["./api"] 