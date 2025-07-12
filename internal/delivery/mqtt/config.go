package mqtt

import (
	"fmt"
	"os"
	"strconv"
)

type MQTTConfig struct {
	BrokerURL string
	ClientID  string
	Username  string
	Password  string
	Port      int
	Topic     string
}

func LoadMqttConfig() *MQTTConfig {
	port, _ := strconv.Atoi(getEnv("MQTT_PORT", "1883"))

	return &MQTTConfig{
		BrokerURL: getEnv("MQTT_BROKER", "localhost"),
		ClientID:  getEnv("MQTT_CLIENT_ID", "vehicle-tracker-client"),
		Username:  getEnv("MQTT_USERNAME", ""),
		Password:  getEnv("MQTT_PASSWORD", ""),
		Port:      port,
		Topic:     getEnv("MQTT_TOPIC", "fleet/vehicle/+/location"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *MQTTConfig) GetBrokerURL() string {
	return fmt.Sprintf("tcp://%s:%d", c.BrokerURL, c.Port)
}
