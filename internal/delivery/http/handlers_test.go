package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/satryo-pramahardi/go-vehicle-tracker/internal/model"
)

type mockVehicleRepo struct {
	mock.Mock
}

func (m *mockVehicleRepo) InsertLocation(loc *model.VehicleLocation) error {
	args := m.Called(loc)
	return args.Error(0)
}

func (m *mockVehicleRepo) GetLatestLocation(vehicleID string) (*model.VehicleLocation, error) {
	args := m.Called(vehicleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.VehicleLocation), args.Error(1)
}

func (m *mockVehicleRepo) GetLocationHistory(vehicleID string, start, end time.Time) ([]*model.VehicleLocation, error) {
	args := m.Called(vehicleID, start, end)
	return nil, args.Error(1)
}

func TestGetLatestLocation_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mockVehicleRepo)
	handler := NewVehicleHandler(mockRepo)

	vehicleID := "TEST123"
	loc := &model.VehicleLocation{
		ID:        1,
		VehicleID: vehicleID,
		Latitude:  -6.2,
		Longitude: 106.8,
		Timestamp: time.Now(),
	}
	mockRepo.On("GetLatestLocation", vehicleID).Return(loc, nil)

	r := gin.Default()
	r.GET("/vehicles/:vehicle_id/location", handler.GetLatestLocation)

	req, _ := http.NewRequest("GET", "/vehicles/TEST123/location", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
	assert.Contains(t, w.Body.String(), vehicleID)
}

func TestGetLatestLocation_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mockVehicleRepo)
	handler := NewVehicleHandler(mockRepo)

	vehicleID := "NOTFOUND"
	mockRepo.On("GetLatestLocation", vehicleID).Return(nil, errors.New("not found"))

	r := gin.Default()
	r.GET("/vehicles/:vehicle_id/location", handler.GetLatestLocation)

	req, _ := http.NewRequest("GET", "/vehicles/NOTFOUND/location", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "vehicle not found")
}
