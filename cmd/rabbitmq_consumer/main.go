package main

import (
	"encoding/json"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type GeofenceAlert struct {
	EventType string  `json:"event_type"`
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	// Connect to RabbitMQ
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://admin:password@localhost:5672/"
	}

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	// Declare the queue
	q, err := ch.QueueDeclare(
		"geofence.event", // queue name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Consume messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Println("🚨 Geofence Alert Consumer Started")
	log.Println("📡 Listening for geofence events on queue: geofence.event")
	log.Println("⏳ Waiting for messages...")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var alert GeofenceAlert
			if err := json.Unmarshal(d.Body, &alert); err != nil {
				log.Printf("❌ Failed to unmarshal alert: %v", err)
				continue
			}

			log.Printf("🚨 GEOFENCE ALERT RECEIVED!")
			log.Printf("   Event Type: %s", alert.EventType)
			log.Printf("   Vehicle ID: %s", alert.VehicleID)
			log.Printf("   Location: (%.4f, %.4f)", alert.Latitude, alert.Longitude)
			log.Printf("   Timestamp: %d", alert.Timestamp)
			log.Printf("   ---")

			// Here you could:
			// - Send SMS notification
			// - Update dashboard
			// - Log to external system
			// - Trigger other business logic
		}
	}()

	<-forever
}
