package repositories

import (
	"dwimc/internal/model"
	testutils "dwimc/internal/test_utils"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
			assertCreatedDevice(t, testDevice, device, err)
		})

		t.Run("multiple", func(t *testing.T) {
			testDevices := generateTestDevices(withNumber(100), withSpecificUser(1))
			for _, testDevice := range testDevices {
				device, err := repo.Create(model.ID(testDevice.userID), testDevice.serial,
					testDevice.name, testDevice.token)
				assertCreatedDevice(t, &testDevice, device, err)
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
				assertCreatedDevice(t, device1, device, err)

				_, err = repo.Create(model.ID(device2.userID), device2.serial,
					device2.name, device2.token)
				assert.Errorf(t, err, "Expected error for duplicate serial")
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
				assertCreatedDevice(t, device1, device, err)

				device, err = repo.Create(model.ID(device2.userID), device2.serial,
					device2.name, device2.token)
				assertCreatedDevice(t, device2, device, err)
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
				assertCreatedDevice(t, device1, device, err)

				_, err = repo.Create(model.ID(device2.userID), device2.serial,
					device2.name, device2.token)

				assert.Errorf(t, err, "Expected error for duplicate token")
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
				assertCreatedDevice(t, device1, device, err)

				_, err = repo.Create(model.ID(device2.userID), device2.serial,
					device2.name, device2.token)

				assert.Errorf(t, err, "Expected error for duplicate token")
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
			assertCreatedDevice(t, &testDevice, device, err)

			retrieved, err := repo.Get(device.ID)

			assert.NoErrorf(t, err, "Get Device failed: %v", err)
			assert.NotNilf(t, retrieved, "Get Device failed - device is nil")

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

			assert.NoErrorf(t, err, "Get Device failed: %v", err)
			assert.Nilf(t, retrieved, "Got: %v, expected: nil", retrieved)
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
				assert.NoErrorf(t, err, "Create Device failed: %v device: %+v", err, device)
			}

			device, err := repo.GetBySerial(model.ID(40), serial)
			assertCreatedDevice(t, &devices[4], device, err)
		})

		t.Run("by serial - none", func(t *testing.T) {
			t.Parallel()

			retrieved, err := repo.GetBySerial(model.ID(100000002), uuid.NewString())

			assert.NoErrorf(t, err, "Get Device failed: %v", err)
			assert.Nilf(t, retrieved, "Got: %v, expected: nil", retrieved)
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
				assert.NoErrorf(t, err, "Create Device failed: %v device: %+v", err, device)
			}

			retrieved, err := repo.GetAllByUserID(model.ID(10))

			assert.NoErrorf(t, err, "Get Devices failed: %v", err)
			assert.Equal(t, 10, len(retrieved), "Expected getting %d devices", 10)
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
				assert.NoErrorf(t, err, "Create Device failed: %v device: %+v", err, device)
			}

			retrieved, err := repo.GetAllByUserID(model.ID(99999))

			assert.NoErrorf(t, err, "Get Devices failed: %v", err)
			assert.Equal(t, 0, len(retrieved), "Expected getting %d devices", 0)
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

			assert.NoErrorf(t, err, "Create Device failed: %v", err)
			assert.NotNil(t, device, "Update Device failed - device is nil")

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

			assert.NoErrorf(t, err, "Update Device failed: %v", err)
			assert.NotNil(t, updated, "Update Device failed - device is nil")

			assert.Equal(t, newSerial, updated.Serial, "Update Device failed - serial not updated")
			assert.Equal(t, newName, updated.Name, "Update Device failed - name not updated")
			assert.Equal(t, newToken, updated.Token.String, "Update Device failed - token not updated")

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

			assert.NoErrorf(t, err, "Create Device failed: %v", err)
			assert.NotNil(t, device, "Create Device failed - device is nil")

			// sets up delay for fooling updated_at field in db.
			time.Sleep(UPDATE_SLEEP_DURATION)

			_, err = repo.Update(device.ID)
			assert.Errorf(t, err, "Expected error for missing fields")
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

			assert.NoErrorf(t, err, "Create Device failed: %v", err)
			assert.NotNil(t, device, "Update Device failed - device is nil")

			err = repo.Delete(device.ID)

			assert.NoErrorf(t, err, "Delete Device failed: %v", err)

			retrieved, err := repo.Get(device.ID)

			assert.NoErrorf(t, err, "Get Device failed: %v", err)
			assert.Nilf(t, retrieved, "Got: %v, expected: nil", retrieved)
		})

		t.Run("by id - none", func(t *testing.T) {
			t.Parallel()

			devices := generateTestDevices(
				withNumber(10),
				withSpecificUser(10),
			)

			for _, device := range devices {
				_, err := repo.Create(model.ID(device.userID), device.serial, device.name, device.token)
				assert.NoErrorf(t, err, "Create Device failed: %v device: %+v", err, device)
			}

			err := repo.Delete(model.ID(99999999))

			assert.NoErrorf(t, err, "Delete Device failed: %v", err)

			retrieved, err := repo.GetAllByUserID(model.ID(10))

			assert.NoErrorf(t, err, "Get Devices failed: %v", err)
			assert.Equal(t, 10, len(retrieved), "Expected getting %d devices", 10)
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
				assert.NoErrorf(t, err, "Create Device failed: %v device: %+v", err, device)
			}

			total, err := repo.DeleteAllByUserID(model.ID(20))

			assert.NoErrorf(t, err, "Delete Devices failed: %v", err)
			assert.Equal(t, 10, int(total), "Expected getting %d devices", 10)

			retrieved, err := repo.GetAllByUserID(model.ID(20))

			assert.NoErrorf(t, err, "Get Devices failed: %v", err)
			assert.Equal(t, 0, len(retrieved), "Expected getting %d devices", 0)
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
				assert.NoErrorf(t, err, "Create Device failed: %v device: %+v", err, device)
			}

			total, err := repo.DeleteAllByUserID(model.ID(999999))

			assert.NoErrorf(t, err, "Delete Devices failed: %v", err)
			assert.Equal(t, 0, int(total), "Expected getting %d devices", 0)

			retrieved, err := repo.GetAllByUserID(model.ID(40))

			assert.NoErrorf(t, err, "Get Devices failed: %v", err)
			assert.Equal(t, 10, len(retrieved), "Expected getting %d devices", 10)
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

func assertCreatedDevice(t *testing.T, expected *testDevice, actual *model.Device, err error) {
	assert.NoErrorf(t, err, "Create Device failed: %v", err)
	assert.NotNilf(t, actual, "Create Device failed - device is nil")

	assert.NotEqualf(t, 0, actual.ID, "Invalid ID 0")

	assert.Equalf(t, model.ID(expected.userID), actual.UserID, "ID mismatch (device ID: %d)", expected.ID)
	assert.Equalf(t, expected.serial, actual.Serial, "Serial mismatch (device ID: %d)", expected.ID)
	assert.Equalf(t, expected.name, actual.Name, "Name mismatch (device ID: %d)", expected.ID)
	assert.Equalf(t, expected.token, actual.Token.String, "Token mismatch (device ID: %d)", expected.ID)
}
