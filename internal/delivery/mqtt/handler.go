package mqtt

import (
	"log"

	"github.com/redis/go-redis/v9"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/app/service"
)

func MessageHandler(rdb *redis.Client) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("[MQTT_HANDLER] Received message on topic %s: %s", msg.Topic(), string(msg.Payload()))
		go func() {
			if err := service.PushLocationUpdateToRedis(rdb, "location_update", "mqtt-subscriber", msg.Payload()); err != nil {
				log.Printf("[MQTT_HANDLER] Failed to push raw event to Redis: %v", err)
			}
		}()
	}
}
