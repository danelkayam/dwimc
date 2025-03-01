package repositories

import (
	"dwimc/internal/model"
	testutils "dwimc/internal/test_utils"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
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
			t.Run("serial - same user", func(t *testing.T) {
				t.Parallel()
				device1 := &testDevice{
					userID: 100001,
					serial: "duplicate-device-serial-100001",
					name:   "device-name-100001-01",
					token:  uuid.NewString(),
				}

				device2 := &testDevice{
					userID: 100001,
					serial: "duplicate-device-serial-100001",
					name:   "device-name-100001-02",
					token:  uuid.NewString(),
				}

				device, err := repo.Create(model.ID(device1.userID), device1.serial,
					device1.name, device1.token)
				checkCreatedDevice(t, device, device1, err)

				_, err = repo.Create(model.ID(device2.userID), device2.serial,
					device2.name, device2.token)

				if err == nil {
					t.Fatalf("Got: nil, expected: error for duplicate fields")
				}
			})

			t.Run("serial - multiple users", func(t *testing.T) {
				t.Parallel()

				sameSerial := uuid.NewString()
				device1 := &testDevice{
					userID: 100002,
					serial: sameSerial,
					name:   "device-name-100002-01",
					token:  uuid.NewString(),
				}

				device2 := &testDevice{
					userID: 100003,
					serial: sameSerial,
					name:   "device-name-100003-02",
					token:  uuid.NewString(),
				}

				device, err := repo.Create(model.ID(device1.userID), device1.serial,
					device1.name, device1.token)
				checkCreatedDevice(t, device, device1, err)

				device, err = repo.Create(model.ID(device2.userID), device2.serial,
					device2.name, device2.token)
				checkCreatedDevice(t, device, device2, err)
			})

			t.Run("token - same user", func(t *testing.T) {
				t.Parallel()

				sameToken := uuid.NewString()
				device1 := &testDevice{
					userID: 100004,
					serial: uuid.NewString(),
					name:   "device-name-100004-01",
					token:  sameToken,
				}

				device2 := &testDevice{
					userID: 100004,
					serial: uuid.NewString(),
					name:   "device-name-100004-02",
					token:  sameToken,
				}

				device, err := repo.Create(model.ID(device1.userID), device1.serial,
					device1.name, device1.token)
				checkCreatedDevice(t, device, device1, err)

				_, err = repo.Create(model.ID(device2.userID), device2.serial,
					device2.name, device2.token)

				if err == nil {
					t.Fatalf("Got: nil, expected: error for duplicate token")
				}
			})

			t.Run("token - multiple users", func(t *testing.T) {
				t.Parallel()

				sameToken := uuid.NewString()
				device1 := &testDevice{
					userID: 100005,
					serial: uuid.NewString(),
					name:   "device-name-100004-01",
					token:  sameToken,
				}

				device2 := &testDevice{
					userID: 100006,
					serial: uuid.NewString(),
					name:   "device-name-100004-02",
					token:  sameToken,
				}

				device, err := repo.Create(model.ID(device1.userID), device1.serial,
					device1.name, device1.token)
				checkCreatedDevice(t, device, device1, err)

				_, err = repo.Create(model.ID(device2.userID), device2.serial,
					device2.name, device2.token)

				if err == nil {
					t.Fatalf("Got: nil, expected: error for duplicate token")
				}
			})
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
			checkCreatedDevice(t, device, &testDevice, err)

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
				t.Fatalf("Get Device failed: %v", err)
			}

			if retrieved != nil {
				t.Fatalf("Got: %v, expected: nil", retrieved)
			}
		})

		t.Run("by serial", func(t *testing.T) {
			t.Parallel()

			devices := generateTestDevices(
				withStartIndex(400),
				withNumber(10),
				withSpecificUser(40),
			)

			// sets specific device with the serial we will look for
			serial := uuid.NewString()
			devices[4].serial = serial

			for _, device := range devices {
				_, err := repo.Create(model.ID(device.userID), device.serial, device.name, device.token)
				if err != nil {
					t.Fatalf("Create Device failed: %v device: %+v", err, device)
				}
			}

			device, err := repo.GetBySerial(model.ID(40), serial)
			checkCreatedDevice(t, device, &devices[4], err)
		})

		t.Run("by serial - none", func(t *testing.T) {
			t.Parallel()

			retrieved, err := repo.GetBySerial(model.ID(100000002), uuid.NewString())
			if err != nil {
				t.Fatalf("Get Device failed: %v", err)
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
					withSpecificUser(10),
				),
				generateTestDevices(
					withStartIndex(200),
					withNumber(10),
					withSpecificUser(20))...,
			)

			for _, device := range devices {
				_, err := repo.Create(model.ID(device.userID), device.serial, device.name, device.token)
				if err != nil {
					t.Fatalf("Create Device failed: %v device: %+v", err, device)
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

			// TODO - fix this - should be simpler

			if updated.Serial != newSerial {
				t.Fatalf("Update Device failed - token not updated")
			}

			if updated.Name != newName {
				t.Fatalf("Update Device failed - name not updated")
			}

			if updated.Token.String != newToken {
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
	number         int
	startIndex     int
	specificUser   bool
	specificUserID uint
}

type testDeviceOption func(options *testDeviceOptions)

func withNumber(number uint) testDeviceOption {
	return func(options *testDeviceOptions) {
		(*options).number = int(number)
	}
}

func withStartIndex(startIndex uint) testDeviceOption {
	return func(options *testDeviceOptions) {
		(*options).startIndex = int(startIndex)
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

	devices := make([]testDevice, options.number)
	startDeviceIndex := options.startIndex

	for i := 0; i < options.number; i++ {
		userID := getUserID(uint(i))
		devices[i] = testDevice{
			ID:     uint(startDeviceIndex + i),
			userID: userID,
			serial: fmt.Sprintf("serial-%d-%d", int(userID), int(i)),
			name:   fmt.Sprintf("name-%d", i),
			token:  fmt.Sprintf("token-%d-%d", int(userID), int(i)),
		}
	}

	return devices
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

	if device.Token.String != testDevice.token {
		t.Fatalf("Got: %s, expected: %s (device ID: %d)",
			device.Token.String, testDevice.token, device.ID)
	}
}
