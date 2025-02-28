package repositories

import (
	"dwimc/internal/model"
	testutils "dwimc/internal/test_utils"
	"fmt"
	"testing"
	"time"
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

			devices := append(
				generateTestDevices(
					withStartIndex(100),
					withNumber(10),
					withSpecificUser(10)),
				generateTestDevices(
					withStartIndex(200),
					withNumber(10),
					withSpecificUser(20))...,
			)

			for _, device := range devices {
				_, err := repo.Create(model.ID(device.userID), device.serial, device.name, device.token)
				if err != nil {
					t.Fatalf("Create Device failed: %v", err)
				}
			}

			retrieved, err := repo.GetAllByUserID(model.ID(10))
			if err != nil {
				t.Fatalf("Get Devices failed: %v", err)
			}

			if len(retrieved) != 10 {
				t.Fatalf("Got: %d, expected: %d", len(retrieved), 10)
			}
		})

		t.Run("all by user id - none", func(t *testing.T) {
			t.Parallel()

			devices := generateTestDevices(
				withStartIndex(300),
				withNumber(10),
				withSpecificUser(30),
			)

			for _, device := range devices {
				_, err := repo.Create(model.ID(device.userID), device.serial, device.name, device.token)
				if err != nil {
					t.Fatalf("Create Device failed: %v", err)
				}
			}

			retrieved, err := repo.GetAllByUserID(model.ID(99999))
			if err != nil {
				t.Fatalf("Get Devices failed: %v", err)
			}

			if len(retrieved) != 0 {
				t.Fatalf("Got: %d, expected: %d", len(retrieved), 0)
			}
		})
	})

	t.Run("update device", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewSQLDeviceRepository(db)

		const UPDATE_SLEEP_DURATION = 1 * time.Second

		testDevices := generateTestDevices(withNumber(10), withSpecificUser(1))

		t.Run("all fields", func(t *testing.T) {
			t.Parallel()

			testDevice := testDevices[0]

			device, err := repo.Create(model.ID(testDevice.userID), testDevice.serial,
				testDevice.name, testDevice.token)

			if err != nil {
				t.Fatalf("Create Device failed: %v", err)
			}

			if device == nil {
				t.Fatalf("Create Device failed - device is nil")
			}

			// sets up delay for fooling updated_at field in db.
			time.Sleep(UPDATE_SLEEP_DURATION)

			newSerial := fmt.Sprintf("new-serial-%d", device.ID)
			newName := fmt.Sprintf("new-name-%d", device.ID)
			newToken := fmt.Sprintf("new-token-%d", device.ID)

			updated, err := repo.Update(device.ID,
				DeviceUpdate.WithSerial(newSerial),
				DeviceUpdate.WithName(newName),
				DeviceUpdate.WithToken(newToken),
			)

			if err != nil {
				t.Fatalf("Update Device failed: %v", err)
			}

			if updated == nil {
				t.Fatalf("Update Device failed - device is nil")
			}

			if updated.Serial != newSerial {
				t.Fatalf("Update Device failed - token not updated")
			}

			if updated.Name != newName {
				t.Fatalf("Update Device failed - name not updated")
			}

			if updated.Token != newToken {
				t.Fatalf("Update Device failed - Token not updated")
			}

			testutils.AssertEqualItems(updated, device,
				func(field string, shouldBeEqual bool, got, expected interface{}) {
					t.Helper()
					if shouldBeEqual {
						t.Fatalf("Mismatch in field %q: got %v, expected %v", field, got, expected)
					} else {
						t.Fatalf("Field %q should have changed, but it did not: got %v", field, got)
					}
				},
				testutils.WithFieldNotEqual[model.Device]("UpdatedAt"),
				testutils.WithFieldNotEqual[model.Device]("Serial"),
				testutils.WithFieldNotEqual[model.Device]("Name"),
				testutils.WithFieldNotEqual[model.Device]("Token"),
			)
		})

		t.Run("none", func(t *testing.T) {
			t.Parallel()

			testDevice := testDevices[1]

			device, err := repo.Create(model.ID(testDevice.userID),
				testDevice.serial, testDevice.name, testDevice.token)

			if err != nil {
				t.Fatalf("Create Device failed: %v", err)
			}

			if device == nil {
				t.Fatalf("Create Device failed - device is nil")
			}

			// sets up delay for fooling updated_at field in db.
			time.Sleep(UPDATE_SLEEP_DURATION)

			_, err = repo.Update(device.ID)
			if err == nil {
				t.Fatalf("Got: nil, expected: error for missing fields")
			}
		})
	})

	t.Run("delete device", func(t *testing.T) {
		db := setupTestDB(t)
		repo := NewSQLDeviceRepository(db)

		t.Run("by id", func(t *testing.T) {
			t.Parallel()

			testDevice := &testDevice{
				userID: 100000,
				serial: "valid-device-serial",
				name:   "valid-device-name",
				token:  "valid-device-token",
			}

			device, err := repo.Create(model.ID(testDevice.userID), testDevice.serial,
				testDevice.name, testDevice.token)

			if err != nil {
				t.Fatalf("Create Device failed: %v", err)
			}

			if device == nil {
				t.Fatalf("Create Device failed - device is nil")
			}

			err = repo.Delete(device.ID)

			if err != nil {
				t.Fatalf("Delete Device failed: %v", err)
			}

			retrieved, err := repo.Get(device.ID)

			if err != nil {
				t.Fatalf("Get Device failed: %v", err)
			}

			if retrieved != nil {
				t.Fatalf("Got %+v, expected: nil", retrieved)
			}
		})

		t.Run("by id - none", func(t *testing.T) {
			t.Parallel()

			devices := generateTestDevices(
				withNumber(10),
				withSpecificUser(10),
			)

			for _, device := range devices {
				_, err := repo.Create(model.ID(device.userID), device.serial, device.name, device.token)
				if err != nil {
					t.Fatalf("Create Device failed: %v", err)
				}
			}

			err := repo.Delete(model.ID(99999999))

			if err != nil {
				t.Fatalf("Delete Device failed: %v", err)
			}

			retrieved, err := repo.GetAllByUserID(model.ID(10))

			if err != nil {
				t.Fatalf("Get Devices failed: %v", err)
			}

			if len(retrieved) != 10 {
				t.Fatalf("Got %d, expected: %d", len(retrieved), 10)
			}
		})

		t.Run("all by userID", func(t *testing.T) {
			t.Parallel()

			devices := append(
				generateTestDevices(
					withStartIndex(200),
					withNumber(10),
					withSpecificUser(20)),
				generateTestDevices(
					withStartIndex(300),
					withNumber(10),
					withSpecificUser(30))...,
			)

			for _, device := range devices {
				_, err := repo.Create(model.ID(device.userID), device.serial, device.name, device.token)
				if err != nil {
					t.Fatalf("Create Device failed: %v", err)
				}
			}

			total, err := repo.DeleteAllByUserID(model.ID(20))

			if err != nil {
				t.Fatalf("Delete Devices failed: %v", err)
			}

			if total != 10 {
				t.Fatalf("Got %+v, expected: %d", total, 10)
			}

			retrieved, err := repo.GetAllByUserID(model.ID(20))

			if err != nil {
				t.Fatalf("Get Devices failed: %v", err)
			}

			if len(retrieved) != 0 {
				t.Fatalf("Got %d, expected: %d", len(retrieved), 0)
			}
		})

		t.Run("all by userID - none", func(t *testing.T) {
			t.Parallel()
			
			devices := generateTestDevices(
				withStartIndex(400),
				withNumber(10),
				withSpecificUser(40),
			)

			for _, device := range devices {
				_, err := repo.Create(model.ID(device.userID), device.serial, device.name, device.token)
				if err != nil {
					t.Fatalf("Create Device failed: %v", err)
				}
			}

			total, err := repo.DeleteAllByUserID(model.ID(999999))

			if err != nil {
				t.Fatalf("Delete Devices failed: %v", err)
			}

			if total != 0 {
				t.Fatalf("Got %+v, expected: %d", total, 0)
			}

			retrieved, err := repo.GetAllByUserID(model.ID(40))

			if err != nil {
				t.Fatalf("Get Devices failed: %v", err)
			}

			if len(retrieved) != 10 {
				t.Fatalf("Got %d, expected: %d", len(retrieved), 10)
			}
		})
	})
}

type testDeviceOptions struct {
	number         uint
	startIndex     uint
	specificUser   bool
	specificUserID uint
}

type testDeviceOption func(options *testDeviceOptions)

func withNumber(number uint) testDeviceOption {
	return func(options *testDeviceOptions) {
		(*options).number = number
	}
}

func withStartIndex(startIndex uint) testDeviceOption {
	return func(options *testDeviceOptions) {
		(*options).startIndex = startIndex
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
		startIndex:     0,
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

	for i := options.startIndex; i < options.number; i++ {
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
