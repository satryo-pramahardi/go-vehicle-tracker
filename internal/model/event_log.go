package model

import (
	"encoding/json"
	"time"
)

type EventEnvelope struct {
	EventType string          `json:"event_type"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
	Source    string          `json:"source"`
}

type EventLog struct {
	ID        int64           `gorm:"primaryKey"`
	EventType string          `gorm:"not null"`
	Timestamp time.Time       `gorm:"index"`
	Payload   json.RawMessage `gorm:"type:jsonb"`
	Source    string
}

type DeadLetterEntry struct {
	EventJSON string `json:"event_json"`
	ErrorMsg  string `json:"error_msg"`
	FailedAt  int64  `json:"failed_at"`
}

func (EventLog) TableName() string {
	return "event_logs"
}
