package main

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/app/service"
)

func main() {
	// Load config from env
	redisAddr := os.Getenv("REDIS_ADDR")
	dbDsn := os.Getenv("DB_DSN")

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Connect to Postgres
	db, err := gorm.Open(postgres.Open(dbDsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Start workers as goroutines
	go service.SaveEventLogFromRedis(rdb, db)
	go service.SaveVehicleLocationFromRedis(rdb, db)
	go service.ArchiveDeadLetterWorker(rdb)

	// Block main from exiting
	select {}
}
