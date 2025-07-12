FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o rabbitmq-consumer cmd/rabbitmq_consumer/main.go

CMD ["./rabbitmq-consumer"] 