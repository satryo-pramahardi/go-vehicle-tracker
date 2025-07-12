package model

import "time"

type VehicleLocation struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	VehicleID string    `gorm:"index" json:"vehicle_id"`
	Latitude  float64   `gorm:"not null" json:"latitude"`
	Longitude float64   `gorm:"not null" json:"longitude"`
	Timestamp time.Time `gorm:"index" json:"timestamp"`
}

func (VehicleLocation) TableName() string {
	return "vehicle_locations"
}
