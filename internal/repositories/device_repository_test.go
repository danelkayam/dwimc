package repositories_test

import (
	"fmt"
	"testing"
)

type testDevice struct {
	ID          uint
	userID      uint
	serial      string
	name        string
	description string
	token       string
}

func TestDeviceRepository(t *testing.T) {

	// TODO - create device
	// TODO - create devices

	// TODO - get device
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
			ID:          uint(i),
			userID:      getUserID(uint(i)),
			serial:      fmt.Sprintf("serial-%d", i),
			name:        fmt.Sprintf("name-%d", i),
			description: fmt.Sprintf("description-%d", i),
			token:       fmt.Sprintf("token-%d", i),
		}
	}

	return testUsers
}
