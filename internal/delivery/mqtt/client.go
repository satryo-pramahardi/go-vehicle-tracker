package mqtt

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client mqtt.Client
	config *MQTTConfig
}

func NewMQTTClient(config *MQTTConfig) *MQTTClient {
	return &MQTTClient{config: config}
}

func (c *MQTTClient) Connect() error {
	options := mqtt.NewClientOptions()
	options.AddBroker(c.config.GetBrokerURL())
	options.SetClientID(c.config.ClientID)

	if c.config.Username != "" {
		options.SetUsername(c.config.Username)
		options.SetPassword(c.config.Password)
	}

	// Add connection timeout and logging
	options.SetConnectTimeout(10 * time.Second)
	options.SetAutoReconnect(true)
	options.SetConnectRetry(true)

	log.Printf("[MQTT CLIENT] Attempting to connect to MQTT broker: %s", c.config.GetBrokerURL())
	c.client = mqtt.NewClient(options)

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		log.Printf("[MQTT CLIENT] MQTT connection failed: %v", token.Error())
		return token.Error()
	}

	log.Printf("[MQTT CLIENT] Successfully connected to MQTT broker")
	log.Printf("[MQTT CLIENT] MQTT client is now listening for messages...")
	return nil
}

func (c *MQTTClient) Subscribe(topic string, callback mqtt.MessageHandler) error {
	log.Printf("[MQTT CLIENT] Subscribing to topic: %s", topic)
	if token := c.client.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
		log.Printf("[MQTT CLIENT] Failed to subscribe: %v", token.Error())
		return token.Error()
	}

	log.Printf("[MQTT CLIENT] Successfully subscribed to topic")
	return nil
}

func (c *MQTTClient) Disconnect() {
	if c.client != nil {
		c.client.Disconnect(250)
		log.Println("[MQTT CLIENT] Disconnected from MQTT broker")
	}
}
