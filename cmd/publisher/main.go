package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

type VehicleLocationPayload struct {
	VehicleID string    `json:"vehicle_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Speed     float64   `json:"speed"`
	Timestamp time.Time `json:"timestamp"`
}

func meterToLatOffset(m float64) float64 {
	return m / 111320 // 1 degree = 111320 meters
}

func main() {

	godotenv.Load()

	brokerEnv := os.Getenv("MQTT_BROKER")
	if brokerEnv == "" {
		brokerEnv = "tcp://mqtt:1883"
	}

	var (
		broker     = flag.String("broker", brokerEnv, "MQTT broker URL")
		vehicleID  = flag.String("vehicle-id", "TJ001", "Vehicle ID")
		interval   = flag.Int("interval", 2, "Seconds between messages")
		count      = flag.Int("count", 0, "Number of messages to send (0=infinite)")
		tripLength = flag.Float64("trip-length", 200, "Trip length in meters before turning")
		speed      = flag.Float64("speed", 5.0, "Vehicle speed")
	)
	flag.Parse()

	log.Printf("[PUBLISHER] Starting publisher, connecting to %s", *broker)
	client := mqtt.NewClient(mqtt.NewClientOptions().AddBroker(*broker))

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("[PUBLISHER] Failed to connect to MQTT broker: %v", token.Error())
	}

	id := *vehicleID
	if id == "" {
		id = fmt.Sprintf("BUS-%03d", rand.Intn(1000))
	}

	baseLat := -6.193125  // Bundaran HI
	baseLon := 106.820233 // Bundaran HI

	topic := os.Getenv("MQTT_TOPIC")
	if topic == "" {
		topic = fmt.Sprintf("fleet/vehicle/%s/location", id)
	}

	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	done := make(chan struct{})
	go func() {
		<-sigChan
		log.Println("[PUBLISHER] Interrupt received, shutting down...")
		done <- struct{}{}
	}()

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	sent := 0
	offset := 0.0
	direction := 1.0
	step := meterToLatOffset(*speed)
	limit := meterToLatOffset(*tripLength)

	for {
		select {
		case <-done:
			client.Disconnect(250)
			log.Printf("Finished sending %d messages", sent)
			return
		case <-ticker.C:
			lat := baseLat + offset
			payload := VehicleLocationPayload{
				VehicleID: id,
				Latitude:  lat,
				Longitude: baseLon,
				Speed:     *speed,
				Timestamp: time.Now(),
			}

			data, err := json.Marshal(payload)
			if err != nil {
				log.Printf("[PUBLISHER] Failed to marshal payload: %v", err)
				continue
			}

			token := client.Publish(topic, 0, false, data)
			token.Wait()

			if token.Error() != nil {
				log.Printf("[PUBLISHER] Failed to publish message: %v", token.Error())
			}

			sent++
			if *count > 0 && sent >= *count {
				client.Disconnect(250)
				log.Printf("Finished sending %d messages", sent)
				return
			}

			offset += step * direction
			if math.Abs(offset) >= limit {
				direction *= -1
			}
			log.Printf("[PUBLISHER] Sent location: vehcile=%s lat=%.6f lon=%.6f speed=%.2f", id, lat, baseLon, *speed)
		}
	}
}
