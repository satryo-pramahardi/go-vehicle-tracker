package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
)

func TestCheckGeofences_InsideGeofence(t *testing.T) {
	// Test vehicle near geofence boundary (100m north of center)
	centerLat := -6.193125
	centerLng := 106.820233
	// 0.0009 degrees latitude ≈ 100m
	vehicleLat := centerLat + 0.0009

	location := model.VehicleLocation{
		VehicleID: "TEST001",
		Latitude:  vehicleLat,
		Longitude: centerLng,
		Timestamp: time.Now(),
	}

	geofence := model.Geofence{
		ID:        1,
		Name:      "Bundaran HI",
		CenterLat: centerLat,
		CenterLng: centerLng,
		Radius:    100,
		Active:    true,
	}

	geofences := []model.Geofence{geofence}

	events := CheckGeofences(location, geofences, nil)

	assert.Len(t, events, 1)
}

func TestCheckGeofences_OutsideGeofence(t *testing.T) {
	location := model.VehicleLocation{
		VehicleID: "TEST002",
		Latitude:  -6.200000,
		Longitude: 106.820233,
		Timestamp: time.Now(),
	}

	geofence := model.Geofence{
		ID:        1,
		Name:      "Bundaran HI",
		CenterLat: -6.193125,
		CenterLng: 106.820233,
		Radius:    100,
		Active:    true,
	}

	geofences := []model.Geofence{geofence}

	events := CheckGeofences(location, geofences, nil)

	assert.Len(t, events, 0)
}

func TestCheckGeofences_AlreadyInside(t *testing.T) {
	location := model.VehicleLocation{
		VehicleID: "TEST003",
		Latitude:  -6.193125,
		Longitude: 106.820233,
		Timestamp: time.Now(),
	}

	geofence := model.Geofence{
		ID:        1,
		Name:      "Bundaran HI",
		CenterLat: -6.193125,
		CenterLng: 106.820233,
		Radius:    100,
		Active:    true,
	}

	geofences := []model.Geofence{geofence}

	events := CheckGeofences(location, geofences, nil)

	assert.Len(t, events, 0)
}

func TestCheckGeofences_NoActiveGeofences(t *testing.T) {
	location := model.VehicleLocation{
		VehicleID: "TEST004",
		Latitude:  -6.193125,
		Longitude: 106.820233,
		Timestamp: time.Now(),
	}

	geofences := []model.Geofence{}

	events := CheckGeofences(location, geofences, nil)

	assert.Len(t, events, 0)
}
