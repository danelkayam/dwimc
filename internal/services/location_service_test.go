package services_test

import (
	"dwimc/internal/model"
	"dwimc/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLocationService(t *testing.T) {
	expectedLocations := generateLocations()
	expected := &expectedLocations[0]

	setupMockService := func() (*MockLocationRepository, services.LocationService) {
		mockRepo := new(MockLocationRepository)
		service := services.NewDefaultLocationService(mockRepo)
		return mockRepo, service
	}

	t.Run("get last location by deviceID", func(t *testing.T) {
		t.Run("valid location", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("GetLastByDeviceID", expected.DeviceID).Return(expected, nil)

			location, err := service.GetLastByDeviceID(expected.DeviceID)
			assert.NoErrorf(t, err, "Get last location failed: %v", err)
			assert.NotNilf(t, location, "Get last location failed - location is nil")
			assert.Equalf(t, expected, location, "Location mismatch")

			mockRepo.AssertExpectations(t)
		})

		t.Run("none", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("GetLastByDeviceID", mock.Anything).Return(nil, model.ErrItemNotFound)

			location, err := service.GetLastByDeviceID(99999)
			assert.ErrorIsf(t, err, model.ErrItemNotFound, "Expected error for item not found")
			assert.Nilf(t, location, "Expected nil location, got: %v", location)
		})
	})

	t.Run("get last location by deviceID", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("GetAllByDeviceID", expected.DeviceID).Return(expectedLocations, nil)

			locations, err := service.GetAllByDeviceID(expected.DeviceID)
			assert.NoErrorf(t, err, "Get locations failed: %v", err)
			assert.Equal(t, 10, len(locations), "Expected getting %d locations", 10)

			mockRepo.AssertExpectations(t)
		})

		t.Run("none", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("GetAllByDeviceID", mock.Anything).Return([]model.Location{}, nil)

			locations, err := service.GetAllByDeviceID(99999)
			assert.NoErrorf(t, err, "Get locations failed: %v", err)
			assert.Equal(t, 0, len(locations), "Expected getting %d locations", 0)
		})
	})

	t.Run("create location", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("Create", expected.DeviceID, expected.Latitude, expected.Longitude).Return(expected, nil)

			location, err := service.Create(expected.DeviceID, expected.Latitude, expected.Longitude)
			assert.NoErrorf(t, err, "Create location failed: %v", err)
			assert.NotNilf(t, location, "Create location failed - location is nil")
			assert.Equalf(t, expected, location, "Location mismatch")

			mockRepo.AssertExpectations(t)
		})

		t.Run("invalid coordinates", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("Create", expected.DeviceID, mock.Anything, mock.Anything).Return(expected, nil)

			location, err := service.Create(expected.DeviceID, 999.99, 999.99) // Invalid lat/lng
			assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected error for invalid coordinates")
			assert.Nilf(t, location, "Expected nil location, got: %v", location)
		})
	})

	t.Run("delete location", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("Delete", expected.ID).Return(nil)

			err := service.Delete(expected.ID)
			assert.NoErrorf(t, err, "Delete location failed: %v", err)

			mockRepo.AssertExpectations(t)
		})

		t.Run("not found", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("Delete", mock.Anything).Return(model.ErrItemNotFound)

			err := service.Delete(99999)
			assert.ErrorIsf(t, err, model.ErrItemNotFound, "Expected error for item not found")
		})
	})

	t.Run("delete all locations by deviceID", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("DeleteAllBy", expected.DeviceID).Return(int64(3), nil)

			deleted, err := service.DeleteAllBy(expected.DeviceID)
			assert.NoErrorf(t, err, "Delete all locations failed: %v", err)
			assert.Equal(t, int64(3), deleted, "Expected deleting %d locations", 3)

			mockRepo.AssertExpectations(t)
		})

		t.Run("none", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("DeleteAllBy", mock.Anything).Return(int64(0), nil)

			deleted, err := service.DeleteAllBy(99999)
			assert.NoErrorf(t, err, "Delete all locations failed: %v", err)
			assert.Equal(t, int64(0), deleted, "Expected deleting %d locations", 0)
		})
	})
}

func generateLocations() []model.Location {
	const size = 10
	locations := []model.Location{}
	date := time.Now()

	for i := 1; i <= size; i++ {
		location := model.Location{
			Model: model.Model{
				ID:        model.ID(i),
				CreatedAt: date,
				UpdatedAt: date,
			},
			DeviceID:  1,
			Latitude:  37.7749 + float64(i)*0.001,
			Longitude: -122.4194 + float64(i)*0.001,
		}

		locations = append(locations, location)
	}

	return locations
}

type MockLocationRepository struct {
	mock.Mock
}

func (m *MockLocationRepository) GetLastByDeviceID(deviceID model.ID) (*model.Location, error) {
	args := m.Called(deviceID)
	if location, ok := args.Get(0).(*model.Location); ok {
		return location, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLocationRepository) GetAllByDeviceID(deviceID model.ID) ([]model.Location, error) {
	args := m.Called(deviceID)
	if locations, ok := args.Get(0).([]model.Location); ok {
		return locations, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLocationRepository) Create(deviceID model.ID, latitude, longitude float64) (*model.Location, error) {
	args := m.Called(deviceID, latitude, longitude)
	if location, ok := args.Get(0).(*model.Location); ok {
		return location, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockLocationRepository) Delete(id model.ID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLocationRepository) DeleteAllBy(deviceID model.ID) (int64, error) {
	args := m.Called(deviceID)
	return args.Get(0).(int64), args.Error(1)
}
