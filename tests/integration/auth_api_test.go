package integration

import (
	"dwimc/internal/model"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthAPI(t *testing.T) {
	const validAPIKey = "8ZZvULIqcPzxwsfnxbWoHUTh"

	router := SetupTestEnv(t, TestEnvParams{
		DatabaseName:         "dwimc_test",
		SecretAPIKey:         validAPIKey,
		LocationHistoryLimit: 10,
	})

	t.Run("valid", func(t *testing.T) {
		devices := PerformOKRequest[[]model.Device](
			t,
			router,
			"GET",
			"/api/devices/",
			validAPIKey,
			nil,
		)

		assert.Equal(t, 0, len(devices), "Missing or too many devices")
	})

	t.Run("unauthenticated", func(t *testing.T) {
		errRes := PerformFailedRequest(
			t,
			router,
			"GET",
			"/api/devices/",
			"blahblah",
			nil,
			http.StatusUnauthorized,
		)

		assert.Equal(t, "Unauthenticated", errRes.Message, "Error message mismatch")
	})
}
