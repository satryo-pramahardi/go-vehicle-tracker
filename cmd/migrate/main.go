package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/db"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	db := db.ConnectGorm()

	log.Println("[MIGRATE] Starting database migration...")

	if err := db.AutoMigrate(&model.VehicleLocation{}); err != nil {
		log.Fatalf("[MIGRATE] Failed to migrate VehicleLocation: %v", err)
	}
	log.Println("[MIGRATE] ✅ VehicleLocation table migrated")

	if err := db.AutoMigrate(&model.EventLog{}); err != nil {
		log.Fatalf("[MIGRATE] Failed to migrate EventLog: %v", err)
	}
	log.Println("[MIGRATE] ✅ EventLog table migrated")

	if err := db.AutoMigrate(&model.Geofence{}); err != nil {
		log.Fatalf("[MIGRATE] Failed to migrate Geofence: %v", err)
	}
	log.Println("[MIGRATE] ✅ Geofence table migrated")

	if err := db.AutoMigrate(&model.GeofenceEvent{}); err != nil {
		log.Fatalf("[MIGRATE] Failed to migrate GeofenceEvent: %v", err)
	}
	log.Println("[MIGRATE] ✅ GeofenceEvent table migrated")

	// Seed some sample geofences
	seedGeofences(db)

	log.Println("[MIGRATE] 🎉 Database migrated successfully")
}

func seedGeofences(db *gorm.DB) {
	// Check if Bundaran HI geofence already exists
	var count int64
	db.Model(&model.Geofence{}).Where("name = ?", "Bundaran HI").Count(&count)
	if count > 0 {
		log.Println("[MIGRATE] 📍 Bundaran HI geofence already exists, skipping seed")
		return
	}

	// Create Bundaran HI geofence
	bundaranHI := model.Geofence{
		Name:      "Bundaran HI",
		CenterLat: -6.2088,
		CenterLng: 106.8456,
		Radius:    100.0, // 100 meters
		Active:    true,
	}

	if err := db.Create(&bundaranHI).Error; err != nil {
		log.Printf("[MIGRATE] ⚠️ Failed to seed Bundaran HI geofence: %v", err)
		return
	}
	log.Println("[MIGRATE] 📍 Bundaran HI geofence seeded successfully")
}
