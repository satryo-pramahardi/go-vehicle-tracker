package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
	"gorm.io/gorm"
)

func SaveEventLogFromRedis(rdb *redis.Client, db *gorm.DB) {
	for {
		envelope, rawEventJSON, err := getEventFromRedis(rdb, "event_log:queue")
		if err != nil {
			log.Printf("[EVENTLOG_WORKER] Error popping event log from Redis: %v", err)
			continue
		}

		eventLog := model.EventLog{
			EventType: envelope.EventType,
			Timestamp: envelope.Timestamp,
			Payload:   envelope.Payload,
			Source:    envelope.Source,
		}

		if err := db.Create(&eventLog).Error; err != nil {
			log.Printf("[EVENTLOG_WORKER] Failed to save event log: %v", err)
			// push to dead letter queue
			pushDeadLetter(rdb, rawEventJSON, err)
		}
	}
}

func pushDeadLetter(rdb *redis.Client, eventJSON string, err error) {
	entry := model.DeadLetterEntry{
		EventJSON: eventJSON,
		ErrorMsg:  err.Error(),
		FailedAt:  time.Now().Unix(),
	}
	entryJSON, _ := json.Marshal(entry)
	_, pushErr := rdb.RPush(context.Background(), "event_log:dead_letter", entryJSON).Result()
	if pushErr != nil {
		log.Printf("[EVENTLOG_WORKER] Error pushing to dead letter queue: %v", pushErr)
	}
}

// Worker that archives dead letter entries to a permanent failed list
func ArchiveDeadLetterWorker(rdb *redis.Client) {
	ctx := context.Background()
	for {
		result, err := rdb.BRPop(ctx, 0, "event_log:dead_letter").Result()
		if err != nil {
			log.Printf("[EVENTLOG_WORKER] Error popping event log from Redis: %v", err)
			continue
		}

		entryJSON := result[1]
		log.Printf("[EVENTLOG_WORKER] Archiving dead letter entry: %s", entryJSON)

		_, err = rdb.RPush(ctx, "event_log:dead_letter_queue", entryJSON).Result()
		if err != nil {
			log.Printf("[EVENTLOG_WORKER] Error pushing to permanent failed list: %v", err)
		}
	}
}
