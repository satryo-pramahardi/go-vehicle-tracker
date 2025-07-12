package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
	"gorm.io/gorm"
)

func SaveVehicleLocationFromRedis(rdb *redis.Client, db *gorm.DB) {
	// Initialize RabbitMQ service
	rabbitMQ, err := NewRabbitMQService("amqp://admin:password@rabbitmq:5672/")
	if err != nil {
		log.Printf("[LOCATION_WORKER] Failed to initialize RabbitMQ: %v", err)
		// Continue without RabbitMQ - geofence events will still be saved to DB
		rabbitMQ = nil
	} else {
		defer rabbitMQ.Close()
		log.Printf("[LOCATION_WORKER] RabbitMQ service initialized successfully")
	}
	for {
		envelope, rawEventJSON, err := getEventFromRedis(rdb, "vehicle_location:queue")
		if err != nil {
			log.Printf("[LOCATION_WORKER] Error popping vehicle location from Redis: %v", err)
			continue
		}

		var vehicleLocation model.VehicleLocation
		err = json.Unmarshal(envelope.Payload, &vehicleLocation)
		if err != nil {
			log.Printf("[LOCATION_WORKER] Failed to unmarshal vehicle location: %v", err)
			errorEnvelope := model.EventEnvelope{
				EventType: "unmarshal_error",
				Source:    "LocationWorker",
				Payload:   []byte(rawEventJSON),
				Timestamp: time.Now(),
			}
			sendEventToRedis(rdb, "event_log:queue", errorEnvelope)
			continue
		}

		log.Printf("[LOCATION_WORKER] Parsed VehicleLocation: %+v", vehicleLocation)

		if err := db.Create(&vehicleLocation).Error; err != nil {
			log.Printf("[LOCATION_WORKER] Failed to save vehicle location: %v", err)
			errorEnvelope := model.EventEnvelope{
				EventType: "save_error",
				Source:    "LocationWorker",
				Payload:   []byte(rawEventJSON),
				Timestamp: time.Now(),
			}
			sendEventToRedis(rdb, "event_log:queue", errorEnvelope)
			continue
		}

		go CallCheckGeofences(vehicleLocation, db, rdb, rabbitMQ)
	}
}
