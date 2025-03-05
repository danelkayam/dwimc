package repositories

import (
	"dwimc/internal/model"
	"testing"

	lom "github.com/samber/lo/mutable"
	"github.com/stretchr/testify/assert"
)

func TestLocationRepo(t *testing.T) {
	t.Run("create location", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewSQLLocationRepository(db)

		t.Run("multiple valid locations", func(t *testing.T) {
			t.Parallel()

			locations := append(
				generateTestLocations(
					generateDevice.withStartIndex(0),
					generateDevice.withNumber(100),
					generateDevice.withSpecificDevice(10),
				),
				generateTestLocations(
					generateDevice.withStartIndex(100),
					generateDevice.withNumber(100),
					generateDevice.withSpecificDevice(20),
				)...,
			)

			for _, location := range locations {
				created, err := repo.Create(location.deviceID, location.latitude, location.longitude)
				assertCreatedLocation(t, &location, created, err)
			}
		})
	})

	t.Run("get location", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewSQLLocationRepository(db)

		t.Run("last", func(t *testing.T) {
			t.Parallel()

			deviceID := model.ID(10)

			locations := generateTestLocations(
				generateDevice.withStartIndex(100),
				generateDevice.withNumber(10),
				generateDevice.withSpecificDevice(deviceID),
			)

			for _, location := range locations {
				created, err := repo.Create(location.deviceID, location.latitude, location.longitude)
				assertCreatedLocation(t, &location, created, err)
			}

			retrieved, err := repo.GetLast(deviceID)
			assert.NoErrorf(t, err, "Get location failed: %v", err)

			assertCreatedLocation(t, &locations[9], retrieved, err)
		})

		t.Run("last - none", func(t *testing.T) {
			t.Parallel()

			locations := generateTestLocations(
				generateDevice.withStartIndex(100),
				generateDevice.withNumber(10),
				generateDevice.withSpecificDevice(20),
			)

			for _, location := range locations {
				created, err := repo.Create(location.deviceID, location.latitude, location.longitude)
				assertCreatedLocation(t, &location, created, err)
			}

			retrieved, err := repo.GetLast(model.ID(999999))
			assert.NoErrorf(t, err, "Get location failed: %v", err)
			assert.Nilf(t, retrieved, "Got: %v, expected: nil", retrieved)
		})

		t.Run("all by device", func(t *testing.T) {
			t.Parallel()

			deviceID := model.ID(30)

			locations := append(
				generateTestLocations(
					generateDevice.withStartIndex(300),
					generateDevice.withNumber(10),
					generateDevice.withSpecificDevice(deviceID),
				),
				generateTestLocations(
					generateDevice.withStartIndex(400),
					generateDevice.withNumber(10),
					generateDevice.withSpecificDevice(40),
				)...,
			)

			for _, location := range locations {
				created, err := repo.Create(location.deviceID, location.latitude, location.longitude)
				assertCreatedLocation(t, &location, created, err)
			}

			testLocations := locations[:10]
			lom.Reverse(testLocations)

			retrieved, err := repo.GetAllBy(deviceID)
			assert.NoErrorf(t, err, "Get location failed: %v", err)
			assert.Equal(t, 10, len(retrieved), "Expected getting %d locations", 10)

			for i, location := range testLocations {
				assertCreatedLocation(t, &location, &retrieved[i], err)
			}
		})

		t.Run("all by device - none", func(t *testing.T) {
			t.Parallel()

			locations := generateTestLocations(
				generateDevice.withStartIndex(200),
				generateDevice.withNumber(10),
				generateDevice.withSpecificDevice(20),
			)

			for _, location := range locations {
				created, err := repo.Create(location.deviceID, location.latitude, location.longitude)
				assertCreatedLocation(t, &location, created, err)
			}

			retrieved, err := repo.GetAllBy(model.ID(999999))
			assert.NoErrorf(t, err, "Get location failed: %v", err)
			assert.Equal(t, 0, len(retrieved), "Expected getting %d devices", 0)
		})
	})

	t.Run("delete location", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewSQLLocationRepository(db)

		t.Run("by ID", func(t *testing.T) {
			t.Parallel()
			
			location := testLocation{
				deviceID: 10,
				latitude: 10.20,
				longitude: 30.40,
			}

			created, err := repo.Create(location.deviceID, location.latitude, location.longitude)
			assertCreatedLocation(t, &location, created, err)

			err = repo.Delete(created.ID)
			assert.NoErrorf(t, err, "Delete location failed: %v", err)

			deleted, err := repo.GetLast(10)
			assert.NoErrorf(t, err, "Get location failed: %v", err)
			assert.Nilf(t, deleted, "Got: %v, expected: nil", deleted)
		})

		t.Run("by ID - none", func(t *testing.T) {
			t.Parallel()

			deviceID := model.ID(20)
			
			locations := generateTestLocations(
				generateDevice.withStartIndex(200),
				generateDevice.withNumber(10),
				generateDevice.withSpecificDevice(deviceID),
			)

			for _, location := range locations {
				created, err := repo.Create(location.deviceID, location.latitude, location.longitude)
				assertCreatedLocation(t, &location, created, err)
			}

			err := repo.Delete(model.ID(99999))
			assert.NoErrorf(t, err, "Delete location failed: %v", err)

			retrieved, err := repo.GetAllBy(deviceID)
			assert.NoErrorf(t, err, "Get location failed: %v", err)
			assert.Equal(t, 10, len(retrieved), "Expected getting %d devices", 10)
		})

		t.Run("all by device", func(t *testing.T) {
			t.Parallel()
			
			deviceID := model.ID(30)
			
			locations := generateTestLocations(
				generateDevice.withStartIndex(300),
				generateDevice.withNumber(10),
				generateDevice.withSpecificDevice(deviceID),
			)

			for _, location := range locations {
				created, err := repo.Create(location.deviceID, location.latitude, location.longitude)
				assertCreatedLocation(t, &location, created, err)
			}

			total, err := repo.DeleteAllBy(deviceID)
			assert.NoErrorf(t, err, "Delete location failed: %v", err)
			assert.Equal(t, 10, int(total), "Expected getting %d devices", 10)

			retrieved, err := repo.GetAllBy(deviceID)
			assert.NoErrorf(t, err, "Get location failed: %v", err)
			assert.Equal(t, 0, len(retrieved), "Expected getting %d devices", 0)
		})

		t.Run("all by device - none", func(t *testing.T) {
			t.Parallel()
			
			deviceID := model.ID(40)
			
			locations := generateTestLocations(
				generateDevice.withStartIndex(400),
				generateDevice.withNumber(10),
				generateDevice.withSpecificDevice(deviceID),
			)

			for _, location := range locations {
				created, err := repo.Create(location.deviceID, location.latitude, location.longitude)
				assertCreatedLocation(t, &location, created, err)
			}

			total, err := repo.DeleteAllBy(model.ID(99999))
			assert.NoErrorf(t, err, "Delete location failed: %v", err)
			assert.Equal(t, 0, int(total), "Expected getting %d devices", 0)

			retrieved, err := repo.GetAllBy(deviceID)
			assert.NoErrorf(t, err, "Get location failed: %v", err)
			assert.Equal(t, 10, len(retrieved), "Expected getting %d devices", 10)
		})
	})
}

type testLocation struct {
	ID        model.ID
	deviceID  model.ID
	latitude  float64
	longitude float64
}

type testLocationOptions struct {
	number           int
	startIndex       int
	specificDevice   bool
	specificDeviceID model.ID
}

type generateLocationOption struct{}

type testLocationOption func(options *testLocationOptions)

func (generateLocationOption) withNumber(number uint) testLocationOption {
	return func(options *testLocationOptions) {
		(*options).number = int(number)
	}
}

func (generateLocationOption) withStartIndex(startIndex uint) testLocationOption {
	return func(options *testLocationOptions) {
		(*options).startIndex = int(startIndex)
	}
}

func (generateLocationOption) withSpecificDevice(deviceID model.ID) testLocationOption {
	return func(options *testLocationOptions) {
		(*options).specificDevice = true
		(*options).specificDeviceID = deviceID
	}
}

var generateDevice generateLocationOption

func generateTestLocations(opts ...testLocationOption) []testLocation {
	options := testLocationOptions{
		number:           10,
		startIndex:       0,
		specificDevice:   false,
		specificDeviceID: 0,
	}

	for _, opt := range opts {
		opt(&options)
	}

	getDeviceID := func(deviceID uint) uint {
		if options.specificDevice {
			return uint(options.specificDeviceID)
		}

		return deviceID
	}

	locations := make([]testLocation, options.number)
	startLocationIndex := options.startIndex

	for i := range options.number {
		deviceID := getDeviceID(uint(i))
		locations[i] = testLocation{
			ID:        model.ID(startLocationIndex + i),
			deviceID:  model.ID(deviceID),
			latitude:  10.10,
			longitude: 10.10,
		}
	}

	return locations
}

func assertCreatedLocation(t *testing.T, expected *testLocation, actual *model.Location, err error) {
	assert.NoErrorf(t, err, "Create Location failed: %v", err)
	assert.NotNilf(t, actual, "Create Location failed - device is nil")

	assert.NotEqualf(t, 0, actual.ID, "Invalid ID 0")

	assert.Equalf(t, model.ID(expected.deviceID), actual.DeviceID, "ID mismatch (location ID: %d)", expected.ID)
	assert.Equalf(t, expected.deviceID, actual.DeviceID, "Device ID mismatch (location ID: %d)", expected.ID)
	assert.Equalf(t, expected.latitude, actual.Latitude, "Latitude mismatch (location ID: %d)", expected.ID)
	assert.Equalf(t, expected.longitude, actual.Longitude, "Longitude mismatch (location ID: %d)", expected.ID)
}
