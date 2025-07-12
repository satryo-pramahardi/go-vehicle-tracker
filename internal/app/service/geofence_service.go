package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/geo"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
	"gorm.io/gorm"
)

type GeofenceEvent struct {
	VehicleID    string
	GeofenceID   int64
	GeofenceName string
	EventType    string // "entering_geofence" or "exiting_geofence"
	Location     model.VehicleLocation
}

// CheckGeofences checks if a vehicle location triggers any geofence events
// Uses a buffer zone (radius ± 5m) to avoid excessive DB queries
func CheckGeofences(loc model.VehicleLocation, geofences []model.Geofence, db *gorm.DB) []GeofenceEvent {
	var events []GeofenceEvent
	for _, geofence := range geofences {
		var bufferMin = geofence.Radius - 5
		var bufferMax = geofence.Radius + 5
		distance := geo.Haversine(loc.Latitude, loc.Longitude, geofence.CenterLat, geofence.CenterLng)
		inside := distance <= geofence.Radius

		if distance >= bufferMin && distance <= bufferMax {
			// Near the boundary, check event log for last state
			lastEventType := getLastGeofenceEventType(loc.VehicleID, geofence.ID, db)
			if inside && lastEventType != model.GeofenceEventEntry {
				events = append(events, GeofenceEvent{
					VehicleID:    loc.VehicleID,
					GeofenceID:   geofence.ID,
					GeofenceName: geofence.Name,
					EventType:    model.GeofenceEventEntry,
					Location:     loc,
				})
			} else if !inside && lastEventType != model.GeofenceEventExit {
				events = append(events, GeofenceEvent{
					VehicleID:    loc.VehicleID,
					GeofenceID:   geofence.ID,
					GeofenceName: geofence.Name,
					EventType:    model.GeofenceEventExit,
					Location:     loc,
				})
			}
			// If lastEventType matches, do nothing (no state change)
		} // else: outside buffer, do nothing (assume state hasn't changed)
	}
	return events
}

// getLastGeofenceEventType queries the geofence events table for the last event type for a specific vehicle and geofence
func getLastGeofenceEventType(vehicleID string, geofenceID int64, db *gorm.DB) string {
	var lastEvent model.GeofenceEvent
	err := db.Where("vehicle_id = ? AND geofence_id = ? AND event_type IN (?, ?)",
		vehicleID, geofenceID, model.GeofenceEventEntry, model.GeofenceEventExit).
		Order("timestamp DESC").
		First(&lastEvent).Error

	if err != nil {
		return "" // No previous event found
	}
	return lastEvent.EventType
}

func CallCheckGeofences(loc model.VehicleLocation, db *gorm.DB, rdb *redis.Client, rabbitMQ *RabbitMQService) {
	var geofences []model.Geofence
	db.Where("active = ?", true).Find(&geofences)

	events := CheckGeofences(loc, geofences, db)
	for _, event := range events {
		// Save geofence event to database
		saveGeofenceEvent(event, db, rdb)

		// Publish RabbitMQ alert for entry events (asynchronously)
		if event.EventType == model.GeofenceEventEntry && rabbitMQ != nil {
			go func(e GeofenceEvent) {
				if err := rabbitMQ.PublishGeofenceAlert(rdb, e.VehicleID, e.Location.Latitude, e.Location.Longitude, e.EventType); err != nil {
					log.Printf("⚠️ Failed to publish RabbitMQ alert: %v", err)
				}
			}(event)
		}
	}
}

func saveGeofenceEvent(event GeofenceEvent, db *gorm.DB, rdb *redis.Client) {
	// Save to dedicated geofence events table
	geofenceEvent := model.GeofenceEvent{
		VehicleID:  event.VehicleID,
		GeofenceID: event.GeofenceID,
		EventType:  event.EventType,
		Timestamp:  time.Now(),
		Latitude:   event.Location.Latitude,
		Longitude:  event.Location.Longitude,
	}

	payload, _ := json.Marshal(event)
	if err := db.Create(&geofenceEvent).Error; err != nil {
		log.Printf("Failed to save geofence event: %v", err)
		pushDeadLetter(rdb, string(payload), err)
		return
	}

	envelope := model.EventEnvelope{
		EventType: event.EventType,
		Source:    "geofence_service",
		Payload:   json.RawMessage(payload),
		Timestamp: time.Now(),
	}
	sendEventToRedis(rdb, "event_log:queue", envelope)
}
