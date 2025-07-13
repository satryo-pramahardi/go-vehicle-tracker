package service

import (
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
)

type GeofenceAlert struct {
	EventType string  `json:"event_type"`
	VehicleID string  `json:"vehicle_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

type RabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQService(amqpURL string) (*RabbitMQService, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare geofence events queue
	_, err = ch.QueueDeclare(
		"geofence.event",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQService{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQService) PublishGeofenceAlert(rdb *redis.Client, vehicleID string, latitude, longitude float64, eventType string) error {
	alert := GeofenceAlert{
		EventType: eventType,
		VehicleID: vehicleID,
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: time.Now().Unix(),
	}

	body, err := json.Marshal(alert)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		"",               // default exchange
		"geofence.event", // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)

	if err != nil {
		log.Printf("[RABBITMQ_SERVICE] ⚠️ Failed to publish geofence alert: %v", err)
		return err
	}

	log.Printf("[RABBITMQ_SERVICE] 📡 Published geofence alert: %s for vehicle %s at (%.4f, %.4f)",
		eventType, vehicleID, latitude, longitude)
	envelope := model.EventEnvelope{
		EventType: eventType,
		Source:    "geofence_service",
		Payload:   json.RawMessage(body),
		Timestamp: time.Now(),
	}
	sendEventToRedis(rdb, "event_log:queue", envelope)
	return nil
}

func (r *RabbitMQService) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
