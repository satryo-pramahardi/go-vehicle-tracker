package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
)

func PushLocationUpdateToRedis(rdb *redis.Client, eventType, source string, payload []byte) error {
	envelope := model.EventEnvelope{
		EventType: eventType,
		Source:    source,
		Payload:   json.RawMessage(payload),
		Timestamp: time.Now(),
	}

	errCh := make(chan error, 2)

	// Push to event_log:queue
	go func() {
		errCh <- sendEventToRedis(rdb, "event_log:queue", envelope)
	}()

	// Push to vehicle_location:queue
	go func() {
		errCh <- sendEventToRedis(rdb, "vehicle_location:queue", envelope)
	}()

	// Wait for both operations to complete
	var finalErr error
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			finalErr = err
		}
	}

	return finalErr
}

func UnmarshalEnvelopePayload[T any](data []byte) (T, error) {
	var envelope model.EventEnvelope
	var result T

	if err := json.Unmarshal(data, &envelope); err != nil {
		return result, err
	}
	if err := json.Unmarshal(envelope.Payload, &result); err != nil {
		return result, err
	}
	return result, nil
}

// helper function to push event to redis
func sendEventToRedis(rdb *redis.Client, queueName string, envelope model.EventEnvelope) error {
	data, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	return rdb.LPush(context.Background(), queueName, data).Err()
}

// helper function to retrieve event from redis queue
func getEventFromRedis(rdb *redis.Client, queueName string) (model.EventEnvelope, string, error) {
	res, err := rdb.BRPop(context.Background(), 0, queueName).Result()
	if err != nil {
		return model.EventEnvelope{}, "", err
	}

	if len(res) < 2 {
		return model.EventEnvelope{}, "", fmt.Errorf("invalid result format")
	}

	rawJSON := res[1] // Keep the raw JSON

	var envelope model.EventEnvelope
	err = json.Unmarshal([]byte(rawJSON), &envelope)
	if err != nil {
		return model.EventEnvelope{}, "", err
	}

	return envelope, rawJSON, nil
}
