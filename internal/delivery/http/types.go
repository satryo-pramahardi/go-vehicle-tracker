package http

import "time"

// Response structures
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Error   string                 `json:"error"`
	Code    string                 `json:"code"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type NotFoundError struct {
	Message   string `json:"message"`
	VehicleID string `json:"vehicle_id"`
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// Request/Response structures
type LocationResponse struct {
	VehicleID string    `json:"vehicle_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
	Speed     float64   `json:"speed,omitempty"`
	Heading   float64   `json:"heading,omitempty"`
}

type LocationHistoryRequest struct {
	VehicleID string `json:"vehicle_id"`
	Start     string `json:"start"`
	End       string `json:"end"`
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Service   string    `json:"service"`
	Timestamp time.Time `json:"timestamp"`
}
