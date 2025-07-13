package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	mqtt_handler "github.com/satryo-pramahardi/go-vehicle-tracker/internal/delivery/mqtt"
)

func main() {
	godotenv.Load()

	config := mqtt_handler.LoadMqttConfig()
	client := mqtt_handler.NewMQTTClient(config)

	if err := client.Connect(); err != nil {
		log.Fatalf("[SUBSCRIBER] Failed to connect to MQTT broker: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Subscribe to MQTT topic
	if err := client.Subscribe(config.Topic, mqtt_handler.MessageHandler(rdb)); err != nil {
		log.Fatalf("[SUBSCRIBER] Failed to subscribe to topic: %v", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	client.Disconnect()
}
