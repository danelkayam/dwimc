package repositories

import (
	"dwimc/internal/model"
	testutils "dwimc/internal/test_utils"
	"fmt"
	"testing"
)

type testDevice struct {
	ID     uint
	userID uint
	serial string
	name   string
	token  string
}

func TestDeviceRepository(t *testing.T) {

	t.Run("create device", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewSQLDeviceRepository(db)

		t.Run("single valid device", func(t *testing.T) {
			testDevice := &testDevice{
				userID: 100000,
				serial: "valid-device-serial",
				name:   "valid-device-name",
				token:  "valid-device-token",
			}

			device, err := repo.Create(model.ID(testDevice.userID), testDevice.serial,
				testDevice.name, testDevice.token)
			checkCreatedDevice(t, device, testDevice, err)
		})

		t.Run("multiple", func(t *testing.T) {
			testDevices := generateTestDevices(withNumber(100), withSpecificUser(1))
			for _, testDevice := range testDevices {
				device, err := repo.Create(model.ID(testDevice.userID), testDevice.serial,
					testDevice.name, testDevice.token)
				checkCreatedDevice(t, device, &testDevice, err)
			}
		})

		t.Run("duplicate fields", func(t *testing.T) {
			t.Fatalf("not implemented!")
		})
	})

	t.Run("get device", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewSQLDeviceRepository(db)

		testDevices := generateTestDevices(withNumber(100), withSpecificUser(1))

		t.Run("by id", func(t *testing.T) {
			t.Parallel()

			testDevice := testDevices[0]

			device, err := repo.Create(model.ID(testDevice.userID), testDevice.serial,
				testDevice.name, testDevice.token)
			if err != nil {
				t.Fatalf("Create Device failed: %v", err)
			}

			if device == nil {
				t.Fatalf("Create Device failed - device is nil")
				return
			}

			retrieved, err := repo.Get(device.ID)
			if err != nil {
				t.Fatalf("Get Device failed: %v", err)
			}

			if retrieved == nil {
				t.Fatalf("Get Device failed - device is nil")
				return
			}

			testutils.AssertEqualItems(device, retrieved,
				func(field string, shouldBeEqual bool, got, expected interface{}) {
					t.Helper()
					if shouldBeEqual {
						t.Fatalf("Mismatch in field %q: got %v, expected %v", field, got, expected)
					} else {
						t.Fatalf("Field %q should have changed, but it did not: got %v", field, got)
					}
				})
		})

		t.Run("by id - none", func(t *testing.T) {
			t.Parallel()

			retrieved, err := repo.Get(model.ID(100000002))
			if err != nil {
				t.Fatalf("Get User failed: %v", err)
			}

			if retrieved != nil {
				t.Fatalf("Got: %v, expected: nil", retrieved)
			}
		})

		t.Run("all by user id", func(t *testing.T) {
			t.Parallel()
			t.Fatalf("not implemented!")
		})

		t.Run("all by user id - none", func(t *testing.T) {
			t.Parallel()
			t.Fatalf("not implemented!")
		})
	})

	t.Run("update device", func(t *testing.T) {
		// db := setupTestDB(t)
		// repo := NewSQLDeviceRepository(db)

		t.Run("all fields", func(t *testing.T) {
			t.Parallel()
			t.Fatalf("not implemented!")
		})

		t.Run("none", func(t *testing.T) {
			t.Parallel()
			t.Fatalf("not implemented!")
		})
	})

	t.Run("delete device", func(t *testing.T) {
		// db := setupTestDB(t)
		// repo := NewSQLDeviceRepository(db)

		t.Run("by id", func(t *testing.T) {
			t.Parallel()
			t.Fatalf("not implemented!")
		})

		t.Run("none", func(t *testing.T) {
			t.Parallel()
			t.Fatalf("not implemented!")
		})

		t.Run("all by userID", func(t *testing.T) {
			t.Parallel()
			t.Fatalf("not implemented!")
		})
	})
}

type testDeviceOptions struct {
	number         int
	specificUser   bool
	specificUserID uint
}

type testDeviceOption func(options *testDeviceOptions)

func withNumber(number uint) testDeviceOption {
	return func(options *testDeviceOptions) {
		(*options).number = int(number)
	}
}

func withSpecificUser(userID uint) testDeviceOption {
	return func(options *testDeviceOptions) {
		(*options).specificUser = true
		(*options).specificUserID = userID
	}
}

func generateTestDevices(opts ...testDeviceOption) []testDevice {
	options := testDeviceOptions{
		number:         10,
		specificUser:   false,
		specificUserID: 0,
	}

	for _, opt := range opts {
		opt(&options)
	}

	getUserID := func(userID uint) uint {
		if options.specificUser {
			return options.specificUserID
		}

		return userID
	}

	testUsers := make([]testDevice, options.number)

	for i := 0; i < len(testUsers); i++ {
		testUsers[i] = testDevice{
			ID:     uint(i),
			userID: getUserID(uint(i)),
			serial: fmt.Sprintf("serial-%d", i),
			name:   fmt.Sprintf("name-%d", i),
			token:  fmt.Sprintf("token-%d", i),
		}
	}

	return testUsers
}

func checkCreatedDevice(t *testing.T, device *model.Device, testDevice *testDevice, err error) {
	if err != nil {
		t.Fatalf("Create Device failed: %v", err)
	}

	if device == nil {
		t.Fatalf("Create Device failed - device is nil")
		return
	}

	if device.ID == 0 {
		t.Fatalf("Got 0, expected: a valid ID")
	}

	if device.UserID != model.ID(testDevice.userID) {
		t.Fatalf("Got: %d, expected: %d (device ID: %d)",
			device.UserID, testDevice.userID, device.ID)
	}

	if device.Serial != testDevice.serial {
		t.Fatalf("Got: %s, expected: %s (device ID: %d)",
			device.Serial, testDevice.serial, device.ID)
	}

	if device.Name != testDevice.name {
		t.Fatalf("Got: %s, expected: %s (device ID: %d)",
			device.Name, testDevice.name, device.ID)
	}

	if device.Token != testDevice.token {
		t.Fatalf("Got: %s, expected: %s (device ID: %d)",
			device.Token, testDevice.token, device.ID)
	}
}
