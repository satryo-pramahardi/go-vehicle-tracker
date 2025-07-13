package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	httphandler "github.com/satryo-pramahardi/go-vehicle-tracker/internal/delivery/http"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/repository"
)

type locationResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type vehicleLocation struct {
	VehicleID string    `json:"vehicle_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

// MockVehicleRepository implements repository.VehicleRepository for testing
type MockVehicleRepository struct {
	mock.Mock
}

func (m *MockVehicleRepository) InsertLocation(loc *model.VehicleLocation) error {
	args := m.Called(loc)
	return args.Error(0)
}

func (m *MockVehicleRepository) GetLatestLocation(vehicleID string) (*model.VehicleLocation, error) {
	args := m.Called(vehicleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.VehicleLocation), args.Error(1)
}

func (m *MockVehicleRepository) GetLocationHistory(vehicleID string, start, end time.Time) ([]*model.VehicleLocation, error) {
	args := m.Called(vehicleID, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.VehicleLocation), args.Error(1)
}

func TestVehicleLocationAPI(t *testing.T) {
	// Test setup
	vehicleID := "API_TEST_001"

	// Create mock repository
	mockRepo := new(MockVehicleRepository)

	// Test successful location retrieval
	t.Run("GetLatestLocation", func(t *testing.T) {
		// Setup mock expectations
		expectedLocation := &model.VehicleLocation{
			ID:        1,
			VehicleID: vehicleID,
			Latitude:  -6.193125,
			Longitude: 106.820233,
			Timestamp: time.Now().UTC().Truncate(time.Second),
		}

		mockRepo.On("GetLatestLocation", vehicleID).Return(expectedLocation, nil)

		// Create handler with mock repository
		handler := httphandler.NewVehicleHandler(mockRepo)

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/api/v1/vehicles/:vehicle_id/location", handler.GetLatestLocation)

		// Make request
		req, _ := http.NewRequest("GET", "/api/v1/vehicles/"+vehicleID+"/location", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)

		var resp locationResp
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "success", resp.Message)

		// Verify API response data
		data, _ := json.Marshal(resp.Data)
		var apiLocation vehicleLocation
		_ = json.Unmarshal(data, &apiLocation)

		assert.Equal(t, vehicleID, apiLocation.VehicleID)
		assert.InDelta(t, expectedLocation.Latitude, apiLocation.Latitude, 0.0001)
		assert.InDelta(t, expectedLocation.Longitude, apiLocation.Longitude, 0.0001)

		// Verify mock was called correctly
		mockRepo.AssertExpectations(t)
		t.Logf("Successfully retrieved location for vehicle %s", vehicleID)
	})

	// Test location not found
	t.Run("GetLatestLocation_NotFound", func(t *testing.T) {
		// Setup mock expectations for not found
		mockRepo.On("GetLatestLocation", "NOTFOUND").Return(nil, repository.ErrVehicleNotFound)

		// Create handler with mock repository
		handler := httphandler.NewVehicleHandler(mockRepo)

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/api/v1/vehicles/:vehicle_id/location", handler.GetLatestLocation)

		// Make request
		req, _ := http.NewRequest("GET", "/api/v1/vehicles/NOTFOUND/location", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusNotFound, w.Code)

		type errorResp struct {
			Error string `json:"error"`
			Code  string `json:"code"`
		}
		var resp errorResp
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "vehicle not found", resp.Error)
		assert.Equal(t, "NOT_FOUND", resp.Code)

		// Verify mock was called correctly
		mockRepo.AssertExpectations(t)
		t.Logf("Correctly handled not found case")
	})
}

func TestVehicleLocationHistory(t *testing.T) {
	// Test setup
	vehicleID := "HISTORY_TEST_001"

	// Create mock repository
	mockRepo := new(MockVehicleRepository)

	// Test location history retrieval
	t.Run("GetLocationHistory", func(t *testing.T) {
		// Setup mock expectations
		baseTime := time.Now().UTC().Truncate(time.Second)
		expectedLocations := []*model.VehicleLocation{
			{
				ID:        1,
				VehicleID: vehicleID,
				Latitude:  -6.193125,
				Longitude: 106.820233,
				Timestamp: baseTime,
			},
			{
				ID:        2,
				VehicleID: vehicleID,
				Latitude:  -6.193200,
				Longitude: 106.820300,
				Timestamp: baseTime.Add(1 * time.Minute),
			},
			{
				ID:        3,
				VehicleID: vehicleID,
				Latitude:  -6.193300,
				Longitude: 106.820400,
				Timestamp: baseTime.Add(2 * time.Minute),
			},
		}

		startTime := baseTime.Add(-30 * time.Second)
		endTime := baseTime.Add(3 * time.Minute)

		mockRepo.On("GetLocationHistory", vehicleID, startTime, endTime).Return(expectedLocations, nil)

		// Create handler with mock repository
		handler := httphandler.NewVehicleHandler(mockRepo)

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/api/v1/vehicles/:vehicle_id/history", handler.GetLocationHistory)

		// Make request
		req, _ := http.NewRequest("GET", "/api/v1/vehicles/"+vehicleID+"/history?start="+startTime.Format(time.RFC3339)+"&end="+endTime.Format(time.RFC3339), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)

		var resp locationResp
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "success", resp.Message)

		// Verify response structure
		data, _ := json.Marshal(resp.Data)
		var historyData map[string]interface{}
		_ = json.Unmarshal(data, &historyData)

		assert.Equal(t, vehicleID, historyData["vehicle_id"])
		assert.Equal(t, float64(3), historyData["count"])

		// Verify mock was called correctly
		mockRepo.AssertExpectations(t)
		t.Logf("Successfully retrieved location history for vehicle %s", vehicleID)
	})
}
