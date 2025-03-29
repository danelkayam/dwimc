package integration

import (
	api_model "dwimc/internal/api/model"
	"dwimc/internal/model"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestDeviceAPI(t *testing.T) {
	const ValidAPIKey = "8ZZvULIqcPzxwsfnxbWoHUTh"

	router := SetupTestEnv(t, TestEnvParams{
		DatabaseName:         "dwimc_test",
		SecretAPIKey:         ValidAPIKey,
		LocationHistoryLimit: 10,
	})

	t.Run("Create Device", func(t *testing.T) {
		payload := api_model.CreateDevice{
			Serial: "device-1-serial",
			Name:   "device 1 name",
		}
		device := PerformOKRequest[model.Device](
			t,
			router,
			"POST",
			"/api/devices/",
			ValidAPIKey,
			payload,
		)

		assert.Greater(t, device.CreatedAt.Unix(), int64(0), "CreatedAt should be valid time")
		assert.Greater(t, device.UpdatedAt.Unix(), int64(0), "CreatedAt should be valid time")
		assert.Equalf(t, payload.Serial, device.Serial, "Serial mismatch")
		assert.Equalf(t, payload.Name, device.Name, "Name mismatch")
	})

	t.Run("Update Device", func(t *testing.T) {
		payload := api_model.CreateDevice{
			Serial: "device-2-serial",
			Name:   "device 2 name",
		}
		device := PerformOKRequest[model.Device](
			t,
			router,
			"POST",
			"/api/devices/",
			ValidAPIKey,
			payload,
		)

		time.Sleep(1 * time.Second)

		updated := PerformOKRequest[model.Device](
			t,
			router,
			"POST",
			"/api/devices/",
			ValidAPIKey,
			api_model.CreateDevice{
				Serial: "device-2-serial",
				Name:   "device 2 name mekmek",
			},
		)

		assert.Equal(t, device.ID, updated.ID, "Must be same device ID")
		assert.Equal(t, device.CreatedAt, updated.CreatedAt, "Must be same device ID")
		assert.Greater(t, updated.UpdatedAt, device.UpdatedAt, "UpdateAt must be greater")
		assert.Equal(t, device.Serial, updated.Serial, "Serial mismatch")
		assert.NotEqual(t, device.Name, updated.Name, "Names are the same")
	})

	t.Run("Delete Device", func(t *testing.T) {
		payload := api_model.CreateDevice{
			Serial: "device-3-serial",
			Name:   "device 3 name",
		}
		device := PerformOKRequest[model.Device](
			t,
			router,
			"POST",
			"/api/devices/",
			ValidAPIKey,
			payload,
		)

		operation := PerformOKRequest[api_model.Operation](
			t,
			router,
			"DELETE",
			fmt.Sprintf("/api/devices/%s", device.ID.Hex()),
			ValidAPIKey,
			nil,
		)

		assert.True(t, operation.Success)

		operation = PerformOKRequest[api_model.Operation](
			t,
			router,
			"DELETE",
			fmt.Sprintf("/api/devices/%s", device.ID.Hex()),
			ValidAPIKey,
			nil,
		)

		assert.False(t, operation.Success)
	})

	t.Run("Get Device", func(t *testing.T) {
		t.Run("valid", func(t *testing.T) {
			payload := api_model.CreateDevice{
				Serial: "device-4-serial",
				Name:   "device 4 name",
			}
			device := PerformOKRequest[model.Device](
				t,
				router,
				"POST",
				"/api/devices/",
				ValidAPIKey,
				payload,
			)

			retrieved := PerformOKRequest[model.Device](
				t,
				router,
				"GET",
				fmt.Sprintf("/api/devices/%s", device.ID.Hex()),
				ValidAPIKey,
				nil,
			)

			assert.Equalf(t, device, retrieved, "Device mismatch")
		})

		t.Run("not found", func(t *testing.T) {
			errRes := PerformFailedRequest(
				t,
				router,
				"GET",
				fmt.Sprintf("/api/devices/%s", bson.NewObjectID().Hex()),
				ValidAPIKey,
				nil,
				http.StatusNotFound,
			)

			assert.Equal(t, "Not found", errRes.Message, "Error message mismatch")
		})
	})

	t.Run("Get Devices", func(t *testing.T) {
		devices := PerformOKRequest[[]model.Device](
			t,
			router,
			"GET",
			"/api/devices/",
			ValidAPIKey,
			nil,
		)

		assert.Equal(t, 3, len(*devices), "Missing or too many devices")
	})
}
