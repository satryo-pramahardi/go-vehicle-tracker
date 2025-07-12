package postgres

import (
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/repository"
	"gorm.io/gorm"
)

type eventLogRepository struct {
	db *gorm.DB
}

func NewEventLogRepository(db *gorm.DB) repository.EventLogRepository {
	return &eventLogRepository{db: db}
}

func (r *eventLogRepository) InsertEvent(evt *model.EventLog) error {
	return r.db.Create(evt).Error
}
