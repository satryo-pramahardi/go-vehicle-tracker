package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
	eventpg "github.com/satryo-pramahardi/go-vehicle-tracker/internal/repository/postgres"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func TestEventLogIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping event log integration test in short mode")
	}

	// Setup database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "admin"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "vehicle_tracker"),
		getEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	eventRepo := eventpg.NewEventLogRepository(db)

	// Test data setup
	testEventType := "vehicle_location_updated"
	testSource := "integration_test"
	testTimestamp := time.Now().UTC().Truncate(time.Second)

	// Create test payload
	payload := map[string]interface{}{
		"vehicle_id": "EVENT_TEST_001",
		"latitude":   -6.193125,
		"longitude":  106.820233,
		"speed":      25.5,
	}
	payloadJSON, _ := json.Marshal(payload)

	// Create event log entry
	eventLog := &model.EventLog{
		EventType: testEventType,
		Timestamp: testTimestamp,
		Payload:   payloadJSON,
		Source:    testSource,
	}

	// Clean up test data
	db.Where("source = ?", testSource).Delete(&model.EventLog{})
	defer db.Where("source = ?", testSource).Delete(&model.EventLog{})

	// Test event insertion
	t.Run("InsertEventLog", func(t *testing.T) {
		err := eventRepo.InsertEvent(eventLog)
		assert.NoError(t, err)
		assert.NotZero(t, eventLog.ID, "Event ID should be set after insertion")

		t.Logf("Successfully inserted event log with ID: %d", eventLog.ID)
	})

	// Test event verification
	t.Run("VerifyEventSaved", func(t *testing.T) {
		var savedEvent model.EventLog
		result := db.Where("id = ?", eventLog.ID).First(&savedEvent)
		assert.NoError(t, result.Error)

		assert.Equal(t, testEventType, savedEvent.EventType)
		assert.Equal(t, testSource, savedEvent.Source)
		assert.Equal(t, testTimestamp.Unix(), savedEvent.Timestamp.Unix())

		// Verify payload data
		var savedPayload map[string]interface{}
		err := json.Unmarshal(savedEvent.Payload, &savedPayload)
		assert.NoError(t, err)
		assert.Equal(t, "EVENT_TEST_001", savedPayload["vehicle_id"])
		assert.InDelta(t, -6.193125, savedPayload["latitude"], 0.0001)
		assert.InDelta(t, 106.820233, savedPayload["longitude"], 0.0001)
		assert.InDelta(t, 25.5, savedPayload["speed"], 0.1)

		t.Logf("Verified event log data: %+v", savedEvent)
	})

	// Test query by event type
	t.Run("QueryEventsByType", func(t *testing.T) {
		var events []model.EventLog
		result := db.Where("event_type = ?", testEventType).Find(&events)
		assert.NoError(t, result.Error)
		assert.GreaterOrEqual(t, len(events), 1, "Should find at least one event of this type")

		// Find our specific event
		found := false
		for _, event := range events {
			if event.ID == eventLog.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find our test event in the query results")

		t.Logf("Found %d events of type '%s'", len(events), testEventType)
	})

	// Test query by source
	t.Run("QueryEventsBySource", func(t *testing.T) {
		var events []model.EventLog
		result := db.Where("source = ?", testSource).Find(&events)
		assert.NoError(t, result.Error)
		assert.GreaterOrEqual(t, len(events), 1, "Should find at least one event from this source")

		t.Logf("Found %d events from source '%s'", len(events), testSource)
	})

	// Test query by timestamp range
	t.Run("QueryEventsByTimestamp", func(t *testing.T) {
		startTime := testTimestamp.Add(-1 * time.Minute)
		endTime := testTimestamp.Add(1 * time.Minute)

		var events []model.EventLog
		result := db.Where("timestamp BETWEEN ? AND ?", startTime, endTime).Find(&events)
		assert.NoError(t, result.Error)
		assert.GreaterOrEqual(t, len(events), 1, "Should find at least one event in the time range")

		t.Logf("Found %d events in timestamp range", len(events))
	})
}

func TestEventLogMultipleEvents(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping multiple events integration test in short mode")
	}

	// Setup database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "admin"),
		getEnv("DB_PASSWORD", "password"),
		getEnv("DB_NAME", "vehicle_tracker"),
		getEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	eventRepo := eventpg.NewEventLogRepository(db)
	testSource := "integration_test_multiple"

	// Clean up test data
	db.Where("source = ?", testSource).Delete(&model.EventLog{})
	defer db.Where("source = ?", testSource).Delete(&model.EventLog{})

	// Insert multiple test events
	eventTypes := []string{"vehicle_created", "location_updated", "geofence_entered"}
	var eventIDs []int64

	for i, eventType := range eventTypes {
		payload := map[string]interface{}{
			"vehicle_id": fmt.Sprintf("MULTI_TEST_%03d", i+1),
			"sequence":   i + 1,
		}
		payloadJSON, _ := json.Marshal(payload)

		eventLog := &model.EventLog{
			EventType: eventType,
			Timestamp: time.Now().UTC().Truncate(time.Second),
			Payload:   payloadJSON,
			Source:    testSource,
		}

		err := eventRepo.InsertEvent(eventLog)
		assert.NoError(t, err)
		eventIDs = append(eventIDs, eventLog.ID)

		// Ensure different timestamps
		time.Sleep(100 * time.Millisecond)
	}

	// Verify all events were inserted
	var count int64
	db.Model(&model.EventLog{}).Where("source = ?", testSource).Count(&count)
	assert.Equal(t, int64(3), count, "Should have inserted exactly 3 events")

	// Query and verify all events
	var events []model.EventLog
	result := db.Where("source = ?", testSource).Order("timestamp ASC").Find(&events)
	assert.NoError(t, result.Error)
	assert.Equal(t, 3, len(events), "Should retrieve exactly 3 events")

	t.Logf("Successfully inserted and retrieved %d events", len(events))
}
