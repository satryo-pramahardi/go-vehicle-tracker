// @title Vehicle Tracker API
// @host localhost:8080

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/satryo-pramahardi/go-vehicle-tracker/docs"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/db"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/delivery/http"
	vehiclepg "github.com/satryo-pramahardi/go-vehicle-tracker/internal/repository/postgres"
)

func main() {
	// Load environment variables
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize db and repository
	db := db.ConnectGorm()
	repo := vehiclepg.NewVehicleLocationRepository(db)

	// Initialize handler and router
	handler := http.NewVehicleHandler(repo)
	router := http.SetupRouter(handler)

	log.Printf("[API_SERVER] Starting API server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("[API_SERVER] Failed to start server: %v", err)
	}
}
