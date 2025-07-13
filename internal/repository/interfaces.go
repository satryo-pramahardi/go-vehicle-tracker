package repository

import (
	"errors"
	"time"

	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
)

// Common repository errors
var (
	ErrVehicleNotFound = errors.New("vehicle not found")
	ErrEventNotFound   = errors.New("event not found")
)

type VehicleRepository interface {
	InsertLocation(loc *model.VehicleLocation) error
	GetLatestLocation(vehicleID string) (*model.VehicleLocation, error)
	GetLocationHistory(vehicleID string, start, end time.Time) ([]*model.VehicleLocation, error)
}

type EventLogRepository interface {
	InsertEvent(evt *model.EventLog) error
}
