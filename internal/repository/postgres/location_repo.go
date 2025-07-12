package postgres

import (
	"time"

	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/repository"
	"gorm.io/gorm"
)

type vehicleLocationRepository struct {
	db *gorm.DB
}

func NewVehicleLocationRepository(db *gorm.DB) repository.VehicleRepository {
	return &vehicleLocationRepository{db: db}
}

func (r *vehicleLocationRepository) InsertLocation(loc *model.VehicleLocation) error {
	return r.db.Create(loc).Error
}

func (r *vehicleLocationRepository) GetLatestLocation(vehicleID string) (*model.VehicleLocation, error) {
	var loc model.VehicleLocation
	err := r.db.Where("vehicle_id = ?", vehicleID).Order("timestamp DESC").First(&loc).Error
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *vehicleLocationRepository) GetLocationHistory(vehicleID string, start, end time.Time) ([]*model.VehicleLocation, error) {
	var history []*model.VehicleLocation
	err := r.db.Where("vehicle_id = ? AND timestamp BETWEEN ? AND ?",
		vehicleID, start, end).Order("timestamp ASC").
		Find(&history).Error
	return history, err
}
