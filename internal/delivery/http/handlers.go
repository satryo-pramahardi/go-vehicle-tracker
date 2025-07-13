package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/repository"
)

// Response helper functions
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func ResponseError(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Error: message,
		Code:  getErrorCode(code),
	})
}

func ResponseNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Error: message,
		Code:  "NOT_FOUND",
	})
}

func ResponseBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: message,
		Code:  "BAD_REQUEST",
	})
}

// Helper function to convert HTTP status code to error code
func getErrorCode(statusCode int) string {
	switch statusCode {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusInternalServerError:
		return "INTERNAL_ERROR"
	default:
		return "UNKNOWN_ERROR"
	}
}

type VehicleHandler struct {
	vehicleRepo repository.VehicleRepository
}

func NewVehicleHandler(vehicleRepo repository.VehicleRepository) *VehicleHandler {
	return &VehicleHandler{
		vehicleRepo: vehicleRepo,
	}
}

// GetLatestLocation godoc
// @Summary      Get latest location
// @Description  Get the latest location for a vehicle
// @Tags         vehicles
// @Param        vehicle_id path string true "Vehicle ID"
// @Success      200  {object}  LocationResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /vehicles/{vehicle_id}/location [get]
func (h *VehicleHandler) GetLatestLocation(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	if vehicleID == "" {
		ResponseBadRequest(c, "vehicle_id is required")
		return
	}

	loc, err := h.vehicleRepo.GetLatestLocation(vehicleID)

	// Error Handlers
	if err != nil {
		ResponseNotFound(c, "vehicle not found")
		return
	}
	if loc == nil {
		ResponseNotFound(c, "vehicle not found")
		return
	}

	ResponseSuccess(c, loc)
}

// GetLocationHistory godoc
// @Summary      Get vehicle location history
// @Description  Get the location history for a vehicle within a time range
// @Tags         vehicles
// @Param        vehicle_id path string true "Vehicle ID"
// @Param        start query string true "Start time (RFC3339)"
// @Param        end query string true "End time (RFC3339)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /vehicles/{vehicle_id}/history [get]
func (h *VehicleHandler) GetLocationHistory(c *gin.Context) {
	vehicleID := c.Param("vehicle_id")
	startStr := c.Query("start")
	endStr := c.Query("end")

	if vehicleID == "" {
		ResponseBadRequest(c, "vehicle_id is required")
		return
	}

	if startStr == "" || endStr == "" {
		ResponseBadRequest(c, "both start and end time parameters are required")
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		ResponseBadRequest(c, "invalid start time format, example: 2023-01-01T00:00:00Z")
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		ResponseBadRequest(c, "invalid end time format, example: 2023-01-01T00:00:00Z")
		return
	}

	if start.After(end) {
		ResponseBadRequest(c, "start time must be before end time")
		return
	}

	history, err := h.vehicleRepo.GetLocationHistory(vehicleID, start, end)
	if err != nil {
		ResponseNotFound(c, "vehicle not found")
		return
	}

	ResponseSuccess(c, gin.H{
		"vehicle_id": vehicleID,
		"start_time": start,
		"end_time":   end,
		"count":      len(history),
		"locations":  history,
	})
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Check if the API is up
// @Tags         health
// @Success      200  {object}  map[string]interface{}
// @Router       /healthz [get]
func (h *VehicleHandler) HealthCheck(c *gin.Context) {
	ResponseSuccess(c, gin.H{
		"status": "ok",
	})
}
