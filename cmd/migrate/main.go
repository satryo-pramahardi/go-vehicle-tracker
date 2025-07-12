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

	log.Println("🔄 Starting database migration...")

	if err := db.AutoMigrate(&model.VehicleLocation{}); err != nil {
		log.Fatalf("💣 Failed to migrate VehicleLocation: %v", err)
	}
	log.Println("✅ VehicleLocation table migrated")

	if err := db.AutoMigrate(&model.EventLog{}); err != nil {
		log.Fatalf("💣 Failed to migrate EventLog: %v", err)
	}
	log.Println("✅ EventLog table migrated")

	if err := db.AutoMigrate(&model.Geofence{}); err != nil {
		log.Fatalf("💣 Failed to migrate Geofence: %v", err)
	}
	log.Println("✅ Geofence table migrated")

	if err := db.AutoMigrate(&model.GeofenceEvent{}); err != nil {
		log.Fatalf("💣 Failed to migrate GeofenceEvent: %v", err)
	}
	log.Println("✅ GeofenceEvent table migrated")

	// Seed Bundaran HI geofence
	seedBundaranHIGeofence(db)

	log.Println("🎉 Database migrated successfully")
}

func seedBundaranHIGeofence(db *gorm.DB) {
	var existingGeofence model.Geofence
	if err := db.Where("name = ?", "Bundaran HI").First(&existingGeofence).Error; err == nil {
		log.Println("📍 Bundaran HI geofence already exists, skipping seed")
		return
	}

	geofence := model.Geofence{
		Name:      "Bundaran HI",
		CenterLat: -6.193125,
		CenterLng: 106.820233,
		Radius:    50.0, // 50 meters
		Active:    true,
	}

	if err := db.Create(&geofence).Error; err != nil {
		log.Printf("⚠️ Failed to seed Bundaran HI geofence: %v", err)
	} else {
		log.Println("📍 Bundaran HI geofence seeded successfully")
	}
}
