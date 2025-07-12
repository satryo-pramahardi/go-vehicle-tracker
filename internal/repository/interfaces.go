package repository

import (
	"time"

	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
)

type VehicleRepository interface {
	InsertLocation(loc *model.VehicleLocation) error
	GetLatestLocation(vehicleID string) (*model.VehicleLocation, error)
	GetLocationHistory(vehicleID string, start, end time.Time) ([]*model.VehicleLocation, error)
}

type EventLogRepository interface {
	InsertEvent(evt *model.EventLog) error
}
