package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	vehiclepg "github.com/satryo-pramahardi/go-vehicle-tracker/internal/repository/postgres"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgres://admin:password@postgres:5432/vehicle_tracker?sslmode=disable"
	port := "8080"

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	repo := vehiclepg.NewVehicleLocationRepository(db)
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/vehicles/:vehicle_id/location", func(c *gin.Context) {
		vehicleID := c.Param("vehicle_id")
		loc, err := repo.GetLatestLocation(vehicleID)
		if err != nil || loc == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
			return
		}
		c.JSON(http.StatusOK, loc)
	})

	r.GET("/vehicles/:vehicle_id/history", func(c *gin.Context) {
		vehicleID := c.Param("vehicle_id")
		startStr := c.Query("start")
		endStr := c.Query("end")
		start, err1 := time.Parse(time.RFC3339, startStr)
		end, err2 := time.Parse(time.RFC3339, endStr)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start or end time"})
			return
		}
		history, err := repo.GetLocationHistory(vehicleID, start, end)
		if err != nil || len(history) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "vehicle not found"})
			return
		}
		c.JSON(http.StatusOK, history)
	})

	r.Run(":" + port)
}
