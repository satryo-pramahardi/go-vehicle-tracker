package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(handler *VehicleHandler) *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/healthz", handler.HealthCheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Apply middleware
	router.Use(corsMiddleware())
	router.Use(loggingMiddleware())

	// API routes
	api := router.Group("/api/v1")
	vehicles := api.Group("/vehicles")
	{
		vehicles.GET("/:vehicle_id/location", handler.GetLatestLocation)
		vehicles.GET("/:vehicle_id/history", handler.GetLocationHistory)
	}

	return router
}
