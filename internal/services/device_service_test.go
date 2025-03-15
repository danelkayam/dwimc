package services_test

import (
	"database/sql"
	"dwimc/internal/model"
	"dwimc/internal/services"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeviceService(t *testing.T) {
	expectedDevices := generateDevices()
	expected := &expectedDevices[0]

	setupMockService := func() (*MockDeviceRepository, services.DeviceService) {
		mockRepo := new(MockDeviceRepository)
		service := services.NewDefaultDeviceService(mockRepo)
		return mockRepo, service
	}

	t.Run("create device", func(t *testing.T) {
		t.Run("valid device", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("Create", expected.UserID, expected.Serial,
				expected.Name, mock.AnythingOfType("string")).Return(expected, nil)

			device, err := service.Create(expected.UserID, expected.Serial, expected.Name)
			assert.NoErrorf(t, err, "Create device failed: %v", err)
			assert.NotNilf(t, device, "Create device failed - device is nil")
			assert.Equalf(t, expected, device, "Device mismatch")
			mockRepo.AssertExpectations(t)
		})

		t.Run("invalid device - duplicated", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("Create", expected.UserID, expected.Serial,
				expected.Name, mock.AnythingOfType("string")).Return(nil, model.ErrItemConflict)

			device, err := service.Create(expected.UserID, expected.Serial, expected.Name)
			assert.ErrorIsf(t, err, model.ErrItemConflict, "Expected error for item conflict")
			assert.Nilf(t, device, "Expected nil device, got: %v", device)
		})

		t.Run("invalid fields", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			// sets "valid" device, invalid fields are validated in the service.
			mockRepo.On("Create", expected.UserID, expected.Serial,
				expected.Name, mock.AnythingOfType("string")).Return(expected, nil)

			device, err := service.Create(expected.UserID, " ", " ")
			assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected error for invalid args")
			assert.Nilf(t, device, "Expected nil device, got: %v", device)
		})
	})

	t.Run("get device", func(t *testing.T) {
		t.Run("get by id", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("GetByID", expected.ID).Return(expected, nil)
			mockRepo.On("GetByID", mock.Anything).Return(nil, model.ErrItemNotFound)

			device, err := service.GetByID(1)
			assert.NoErrorf(t, err, "Get device failed: %v", err)
			assert.NotNilf(t, device, "Get device failed - device is nil")
			assert.Equalf(t, device, expected, "device mismatch")
			mockRepo.AssertExpectations(t)
		})

		t.Run("get by id - none", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("GetByID", expected.ID).Return(expected, nil)
			mockRepo.On("GetByID", mock.Anything).Return(nil, model.ErrItemNotFound)

			device, err := service.GetByID(999999)
			assert.ErrorIsf(t, err, model.ErrItemNotFound, "Expected error for item not found")
			assert.Nilf(t, device, "Expected nil device, got: %v", device)
		})

		t.Run("get by serial", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("GetBySerial", expected.UserID, expected.Serial).Return(expected, nil)
			mockRepo.On("GetBySerial", mock.Anything, mock.Anything).Return(nil, model.ErrItemNotFound)

			device, err := service.GetBySerial(1, "device-serial-1-1")
			assert.NoErrorf(t, err, "Get device failed: %v", err)
			assert.NotNilf(t, device, "Get device failed - device is nil")
			assert.Equalf(t, device, expected, "device mismatch")
			mockRepo.AssertExpectations(t)
		})

		t.Run("get by serial - none", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("GetBySerial", expected.UserID, expected.Serial).Return(expected, nil)
			mockRepo.On("GetBySerial", mock.Anything, mock.Anything).Return(nil, model.ErrItemNotFound)

			device, err := service.GetBySerial(99999, "non-existing-serial")
			assert.ErrorIsf(t, err, model.ErrItemNotFound, "Expected error for item not found")
			assert.Nilf(t, device, "Expected nil device, got: %v", device)

			device, err = service.GetBySerial(1, "non-existing-serial")
			assert.ErrorIsf(t, err, model.ErrItemNotFound, "Expected error for item not found")
			assert.Nilf(t, device, "Expected nil device, got: %v", device)

			device, err = service.GetBySerial(99999, expected.Serial)
			assert.ErrorIsf(t, err, model.ErrItemNotFound, "Expected error for item not found")
			assert.Nilf(t, device, "Expected nil device, got: %v", device)
		})
	})

	t.Run("get devices", func(t *testing.T) {
		t.Run("get all by by userID", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("GetAllByUserID", expected.UserID).Return(expectedDevices, nil)

			devices, err := service.GetAllByUserID(expected.UserID)
			assert.NoErrorf(t, err, "Get devices failed: %v", err)
			assert.Equal(t, 10, len(devices), "Expected getting %d devices", 10)
		})

		t.Run("get all by by userID - none", func(t *testing.T) {
			t.Parallel()

			mockRepo, service := setupMockService()
			mockRepo.On("GetAllByUserID", model.ID(99999)).Return([]model.Device{}, nil)

			devices, err := service.GetAllByUserID(99999)
			assert.NoErrorf(t, err, "Get devices failed: %v", err)
			assert.Equal(t, 0, len(devices), "Expected getting %d devices", 0)
		})
	})

	t.Run("update device", func(t *testing.T) {
		t.Run("all fields", func(t *testing.T) {
			t.Parallel()

			expected := &model.Device{
				Model: model.Model{
					ID:        expected.ID,
					CreatedAt: expected.CreatedAt,
					UpdatedAt: expected.UpdatedAt,
				},
				UserID: expected.UserID,
				Serial: "some-serial",
				Name:   "some-name",
				Token:  sql.NullString{String: "some-token-not-care", Valid: true},
			}

			fields := []model.Field{
				model.WithSerial(expected.Serial),
				model.WithName(expected.Name),
			}

			mockRepo, service := setupMockService()
			mockRepo.On("Update", expected.ID, fields).Return(expected, nil)

			device, err := service.Update(expected.ID, fields...)
			assert.NoErrorf(t, err, "Update device failed: %v", err)
			assert.NotNilf(t, device, "Update device failed - device is nil")
			assert.Equalf(t, device, expected, "device mismatch")
			mockRepo.AssertExpectations(t)
		})

		t.Run("partial fields", func(t *testing.T) {
			t.Parallel()

			expected := &model.Device{
				Model: model.Model{
					ID:        expected.ID,
					CreatedAt: expected.CreatedAt,
					UpdatedAt: expected.UpdatedAt,
				},
				UserID: expected.UserID,
				Serial: "some-serial",
				Name:   expected.Name,
				Token:  sql.NullString{String: "some-token-not-care", Valid: true},
			}

			fields := []model.Field{
				model.WithSerial(expected.Serial),
			}

			mockRepo, service := setupMockService()
			mockRepo.On("Update", expected.ID, fields).Return(expected, nil)

			device, err := service.Update(expected.ID, fields...)
			assert.NoErrorf(t, err, "Update device failed: %v", err)
			assert.NotNilf(t, device, "Update device failed - device is nil")
			assert.Equalf(t, device, expected, "device mismatch")
			mockRepo.AssertExpectations(t)
		})

		t.Run("no fields", func(t *testing.T) {
			t.Parallel()

			fields := []model.Field{}

			mockRepo, service := setupMockService()
			mockRepo.On("Update", expected.ID, fields).Return(expected, nil)

			device, err := service.Update(expected.ID, fields...)
			assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected error for invalid fields")
			assert.Nilf(t, device, "Expected nil device, got: %v", device)
		})

		t.Run("invalid field name", func(t *testing.T) {
			t.Parallel()

			fields := []model.Field{
				model.WithField("not_existing_field", "some_value"),
			}

			mockRepo, service := setupMockService()
			mockRepo.On("Update", expected.ID, fields).Return(expected, nil)

			device, err := service.Update(expected.ID, fields...)
			assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected error for invalid fields")
			assert.Nilf(t, device, "Expected nil device, got: %v", device)
		})

		t.Run("invalid fields value", func(t *testing.T) {
			t.Parallel()

			fields := []model.Field{
				model.WithSerial(""),
				model.WithName(""),
			}

			mockRepo, service := setupMockService()
			mockRepo.On("Update", expected.ID, fields).Return(expected, nil)

			device, err := service.Update(expected.ID, fields...)
			assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected error for invalid fields")
			assert.Nilf(t, device, "Expected nil device, got: %v", device)
		})

		t.Run("device not found", func(t *testing.T) {
			t.Parallel()

			fields := []model.Field{
				model.WithSerial(expected.Serial),
				model.WithName(expected.Name),
			}

			mockRepo, service := setupMockService()
			mockRepo.On("Update", model.ID(99999), fields).Return(nil, model.ErrItemNotFound)

			device, err := service.Update(model.ID(99999), fields...)
			assert.ErrorIsf(t, err, model.ErrItemNotFound, "Expected error for item not found")
			assert.Nilf(t, device, "Expected nil device, got: %v", device)
		})

		t.Run("update unchanged fields", func(t *testing.T) {
			t.Parallel()

			fields := []model.Field{
				model.WithSerial(expected.Serial),
				model.WithName(expected.Name),
			}

			mockRepo, service := setupMockService()
			mockRepo.On("Update", expected.ID, fields).Return(expected, nil)

			device, err := service.Update(expected.ID, fields...)
			assert.NoErrorf(t, err, "Update device failed: %v", err)
			assert.NotNilf(t, device, "Update device failed - device is nil")
			assert.Equalf(t, device, expected, "device mismatch")
			mockRepo.AssertExpectations(t)
		})

		t.Run("database error", func(t *testing.T) {
			t.Parallel()

			fields := []model.Field{
				model.WithSerial(expected.Serial),
				model.WithName(expected.Name),
			}

			mockRepo, service := setupMockService()
			mockRepo.On("Update", expected.ID, fields).Return(nil, model.ErrDatabase)

			device, err := service.Update(expected.ID, fields...)
			assert.ErrorIsf(t, err, model.ErrDatabase, "Expected error for database")
			assert.Nilf(t, device, "Expected nil device, got: %v", device)
		})
	})

	t.Run("delete device", func(t *testing.T) {
		t.Run("by id", func(t *testing.T) {
			t.Parallel()

			id := model.ID(1)

			mockRepo, service := setupMockService()
			mockRepo.On("Delete", id).Return(nil)

			err := service.Delete(id)
			assert.NoErrorf(t, err, "Delete device failed: %v", err)
			mockRepo.AssertExpectations(t)
		})

		t.Run("by id - none", func(t *testing.T) {
			t.Parallel()

			id := model.ID(99999)

			mockRepo, service := setupMockService()
			mockRepo.On("Delete", id).Return(model.ErrItemNotFound)

			err := service.Delete(id)
			assert.ErrorIsf(t, err, model.ErrItemNotFound, "Expected error for item not found")
			mockRepo.AssertExpectations(t)
		})

		t.Run("all by userID", func(t *testing.T) {
			t.Parallel()

			userID := model.ID(1)

			mockRepo, service := setupMockService()
			mockRepo.On("DeleteAllByUserID", userID).Return(int64(3), nil)

			deleted, err := service.DeleteAllByUserID(userID)
			assert.NoErrorf(t, err, "Delete device failed: %v", err)
			assert.Equal(t, int64(3), deleted, "Expected deleting %d devices", 3)
			mockRepo.AssertExpectations(t)
		})

		t.Run("all by userID - none", func(t *testing.T) {
			t.Parallel()

			userID := model.ID(99999)

			mockRepo, service := setupMockService()
			mockRepo.On("DeleteAllByUserID", userID).Return(int64(0), nil)

			deleted, err := service.DeleteAllByUserID(userID)
			assert.NoErrorf(t, err, "Delete Devices failed: %v", err)
			assert.Equal(t, int64(0), deleted, "Expected getting %d devices", 0)
			mockRepo.AssertExpectations(t)
		})
	})
}

func generateDevices() []model.Device {
	const size = 10
	devices := []model.Device{}
	date := time.Now()

	for i := 1; i <= size; i++ {
		device := model.Device{
			Model: model.Model{
				ID:        model.ID(i),
				CreatedAt: date,
				UpdatedAt: date,
			},
			UserID: 1,
			Serial: fmt.Sprintf("device-serial-%d-%d", i, i),
			Name:   fmt.Sprintf("device-name-%d-%d", i, i),
			Token:  sql.NullString{String: fmt.Sprintf("device-token-%d-%d", i, i), Valid: true},
		}

		devices = append(devices, device)
	}

	return devices
}

type MockDeviceRepository struct {
	mock.Mock
}

func (m *MockDeviceRepository) GetByID(id model.ID) (*model.Device, error) {
	args := m.Called(id)
	if device, ok := args.Get(0).(*model.Device); ok {
		return device, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDeviceRepository) GetBySerial(userID model.ID, serial string) (*model.Device, error) {
	args := m.Called(userID, serial)
	if device, ok := args.Get(0).(*model.Device); ok {
		return device, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDeviceRepository) GetAllByUserID(userID model.ID) ([]model.Device, error) {
	args := m.Called(userID)
	if devices, ok := args.Get(0).([]model.Device); ok {
		return devices, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDeviceRepository) Create(userID model.ID, serial, name, token string) (*model.Device, error) {
	args := m.Called(userID, serial, name, token)
	if device, ok := args.Get(0).(*model.Device); ok {
		return device, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDeviceRepository) Update(id model.ID, fields ...model.Field) (*model.Device, error) {
	args := m.Called(id, fields)
	if device, ok := args.Get(0).(*model.Device); ok {
		return device, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockDeviceRepository) Delete(id model.ID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockDeviceRepository) DeleteAllByUserID(userID model.ID) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}
