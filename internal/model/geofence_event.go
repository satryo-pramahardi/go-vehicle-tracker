package model

import "time"

// Geofence event types
const (
	GeofenceEventEntry = "geofence_entry"
	GeofenceEventExit  = "geofence_exit"
)

type GeofenceEvent struct {
	ID         int64     `gorm:"primaryKey"`
	VehicleID  string    `gorm:"not null;index:idx_vehicle_geofence"`
	GeofenceID int64     `gorm:"not null;index:idx_vehicle_geofence"`
	EventType  string    `gorm:"not null;check:event_type IN ('geofence_entry', 'geofence_exit')"`
	Timestamp  time.Time `gorm:"not null;index"`
	Latitude   float64   `gorm:"not null"`
	Longitude  float64   `gorm:"not null"`
}

func (GeofenceEvent) TableName() string {
	return "geofence_events"
}
