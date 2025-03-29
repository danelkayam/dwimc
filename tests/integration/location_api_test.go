package integration

import (
	api_model "dwimc/internal/api/model"
	"dwimc/internal/model"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestLocationAPI(t *testing.T) {
	const validAPIKey = "8ZZvULIqcPzxwsfnxbWoHUTh"
	const locationHistory = 5

	router := SetupTestEnv(t, TestEnvParams{
		DatabaseName:         "dwimc_test",
		SecretAPIKey:         validAPIKey,
		LocationHistoryLimit: locationHistory,
	})

	createDevice := func(serial string, name string) model.Device {
		return PerformOKRequest[model.Device](
			t,
			router,
			"POST",
			"/api/devices/",
			validAPIKey,
			api_model.CreateDevice{
				Serial: serial,
				Name:   name,
			},
		)
	}

	createLocation := func(deviceID string, payload api_model.CreateLocation) api_model.Operation {
		return PerformOKRequest[api_model.Operation](
			t,
			router,
			"POST",
			fmt.Sprintf("/api/devices/%s/locations/", deviceID),
			validAPIKey,
			payload,
		)
	}

	t.Run("Create Location", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			device := createDevice("device-1-serial", "device-1-name")

			operation := createLocation(
				device.ID.Hex(),
				api_model.CreateLocation{
					Latitude:  32.086880,
					Longitude: 34.775759,
				},
			)

			assert.True(t, operation.Success)
		})

		t.Run("invalid params", func(t *testing.T) {
			payloads := []api_model.CreateLocation{
				{},
				{
					Latitude:  -200.000000,
					Longitude: 34.775759,
				},
				{
					Latitude:  32.086880,
					Longitude: 200.129387,
				},
				{
					Latitude: 32.086880,
				},
				{
					Longitude: 34.775759,
				},
			}

			device := createDevice("device-1-serial", "device-1-name")

			for _, payload := range payloads {
				errRes := PerformFailedRequest(
					t,
					router,
					"POST",
					fmt.Sprintf("/api/devices/%s/locations/", device.ID.Hex()),
					validAPIKey,
					payload,
					http.StatusBadRequest,
				)

				assert.Equal(t, "Bad request", errRes.Message, "Error message mismatch")
			}
		})

		t.Run("invalid no device", func(t *testing.T) {
			errRes := PerformFailedRequest(
				t,
				router,
				"POST",
				fmt.Sprintf("/api/devices/%s/locations/", bson.NewObjectID().Hex()),
				validAPIKey,
				api_model.CreateLocation{
					Latitude:  32.086880,
					Longitude: 34.775759,
				},
				http.StatusNotFound,
			)

			assert.Equal(t, "Not found", errRes.Message, "Error message mismatch")
		})
	})

	t.Run("Get Last Location", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			payload := api_model.CreateLocation{
				Latitude:  32.086880,
				Longitude: 34.775759,
			}

			device := createDevice("device-2-serial", "device-2-name")

			operation := createLocation(
				device.ID.Hex(),
				payload,
			)

			assert.True(t, operation.Success)

			location := PerformOKRequest[model.Location](
				t,
				router,
				"GET",
				fmt.Sprintf("/api/devices/%s/locations/latest", device.ID.Hex()),
				validAPIKey,
				nil,
			)

			assert.Equal(t, device.ID.Hex(), location.DeviceID.Hex(), "ID mismatch")
			assert.Greater(t, location.CreatedAt.Unix(), int64(0), "CreatedAt should be valid time")
			assert.Greater(t, location.UpdatedAt.Unix(), int64(0), "UpdatedAt should be valid time")
			assert.Equal(t, payload.Latitude, location.Latitude, "Latitude mismatch")
			assert.Equal(t, payload.Longitude, location.Longitude, "Longitude mismatch")
		})

		t.Run("nothing", func(t *testing.T) {
			device := createDevice("device-3-serial", "device-3-name")

			response := PerformOKRequestNoValidateResponse(
				t,
				router,
				"GET",
				fmt.Sprintf("/api/devices/%s/locations/latest", device.ID.Hex()),
				validAPIKey,
				nil,
			)

			assert.Nil(t, response.Data, "Location is not nil")
			assert.Nil(t, response.Error, "Error is not nil")
		})
	})

	t.Run("Get Locations", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			count := 3
			payload := api_model.CreateLocation{
				Latitude:  32.086880,
				Longitude: 34.775759,
			}

			device := createDevice("device-4-serial", "device-4-name")

			for range count {
				operation := createLocation(
					device.ID.Hex(),
					payload,
				)

				assert.True(t, operation.Success)
			}

			locations := PerformOKRequest[[]model.Location](
				t,
				router,
				"GET",
				fmt.Sprintf("/api/devices/%s/locations/", device.ID.Hex()),
				validAPIKey,
				nil,
			)

			assert.Equal(t, count, len(locations))

			for _, location := range locations {
				assert.Equal(t, device.ID.Hex(), location.DeviceID.Hex(), "ID mismatch")
				assert.Greater(t, location.CreatedAt.Unix(), int64(0), "CreatedAt should be valid time")
				assert.Greater(t, location.UpdatedAt.Unix(), int64(0), "UpdatedAt should be valid time")
				assert.Equal(t, payload.Latitude, location.Latitude, "Latitude mismatch")
				assert.Equal(t, payload.Longitude, location.Longitude, "Longitude mismatch")
			}
		})

		t.Run("history limit", func(t *testing.T) {
			count := 3 * locationHistory
			device := createDevice("device-5-serial", "device-5-name")

			for range count {
				operation := createLocation(
					device.ID.Hex(),
					api_model.CreateLocation{
						Latitude:  32.086880,
						Longitude: 34.775759,
					},
				)

				assert.True(t, operation.Success)
			}

			locations := PerformOKRequest[[]model.Location](
				t,
				router,
				"GET",
				fmt.Sprintf("/api/devices/%s/locations/", device.ID.Hex()),
				validAPIKey,
				nil,
			)

			assert.Equal(t, locationHistory, len(locations))
		})
	})

	t.Run("Delete Locations", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			count := 3
			device := createDevice("device-5-serial", "device-5-name")

			for range count {
				operation := createLocation(
					device.ID.Hex(),
					api_model.CreateLocation{
						Latitude:  32.086880,
						Longitude: 34.775759,
					},
				)

				assert.True(t, operation.Success)
			}

			operation := PerformOKRequest[api_model.Operation](
				t,
				router,
				"DELETE",
				fmt.Sprintf("/api/devices/%s/locations/", device.ID.Hex()),
				validAPIKey,
				nil,
			)

			assert.True(t, operation.Success)

			locations := PerformOKRequest[[]model.Location](
				t,
				router,
				"GET",
				fmt.Sprintf("/api/devices/%s/locations/", device.ID.Hex()),
				validAPIKey,
				nil,
			)

			assert.Equal(t, 0, len(locations))
		})

		t.Run("empty", func(t *testing.T) {
			device := createDevice("device-6-serial", "device-6-name")

			operation := PerformOKRequest[api_model.Operation](
				t,
				router,
				"DELETE",
				fmt.Sprintf("/api/devices/%s/locations/", device.ID.Hex()),
				validAPIKey,
				nil,
			)

			assert.False(t, operation.Success)

			locations := PerformOKRequest[[]model.Location](
				t,
				router,
				"GET",
				fmt.Sprintf("/api/devices/%s/locations/", device.ID.Hex()),
				validAPIKey,
				nil,
			)

			assert.Equal(t, 0, len(locations))
		})

		t.Run("delete device", func(t *testing.T) {
			count := 3
			device := createDevice("device-7-serial", "device-7-name")

			for range count {
				operation := createLocation(
					device.ID.Hex(),
					api_model.CreateLocation{
						Latitude:  32.086880,
						Longitude: 34.775759,
					},
				)

				assert.True(t, operation.Success)
			}

			operation := PerformOKRequest[api_model.Operation](
				t,
				router,
				"DELETE",
				fmt.Sprintf("/api/devices/%s", device.ID.Hex()),
				validAPIKey,
				nil,
			)

			assert.True(t, operation.Success)

			errRes := PerformFailedRequest(
				t,
				router,
				"GET",
				fmt.Sprintf("/api/devices/%s/locations/", device.ID.Hex()),
				validAPIKey,
				nil,
				http.StatusNotFound,
			)

			assert.Equal(t, "Not found", errRes.Message, "Error message mismatch")
		})
	})

	t.Run("Delete Location", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			device := createDevice("device-8-serial", "device-8-name")

			operation := createLocation(
				device.ID.Hex(),
				api_model.CreateLocation{
					Latitude:  32.086880,
					Longitude: 34.775759,
				},
			)

			assert.True(t, operation.Success)

			location := PerformOKRequest[model.Location](
				t,
				router,
				"GET",
				fmt.Sprintf("/api/devices/%s/locations/latest", device.ID.Hex()),
				validAPIKey,
				nil,
			)

			operation = PerformOKRequest[api_model.Operation](
				t,
				router,
				"DELETE",
				fmt.Sprintf("/api/devices/%s/locations/%s",
					device.ID.Hex(),
					location.ID.Hex(),
				),
				validAPIKey,
				nil,
			)

			assert.True(t, operation.Success)

			response := PerformOKRequestNoValidateResponse(
				t,
				router,
				"GET",
				fmt.Sprintf("/api/devices/%s/locations/latest", device.ID.Hex()),
				validAPIKey,
				nil,
			)

			assert.Nil(t, response.Data, "Location is not nil")
			assert.Nil(t, response.Error, "Error is not nil")
		})
	})
}
