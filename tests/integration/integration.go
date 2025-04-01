package integration

import (
	"context"
	"dwimc/internal/api"
	"dwimc/internal/database"
	"dwimc/internal/repositories"
	"dwimc/internal/services"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
)

type TestEnvParams struct {
	DatabaseName         string
	SecretAPIKey         string
	LocationHistoryLimit int
}

func SetupTestEnv(t *testing.T, params TestEnvParams) *gin.Engine {
	ctx := context.Background()

	container, err := mongodb.Run(ctx, "mongo:latest")
	require.NoError(t, err, "Failed to create mongodb container")

	uri, err := container.ConnectionString(ctx)
	require.NoError(t, err, "Failed to get mongodb container URI")

	client, err := database.InitializeDatabase(uri)
	require.NoError(t, err, "Failed to initialize mongodb client")

	deviceRepo, err := repositories.NewMongodbDeviceRepository(
		ctx,
		client,
		params.DatabaseName,
	)
	require.NoError(t, err, "Failed to create device repository")

	locationRepo, err := repositories.NewMongodbLocationRepository(
		ctx,
		client,
		params.DatabaseName,
	)
	require.NoError(t, err, "Failed to create location repository")

	router := api.InitializeRouters(
		false,
		params.SecretAPIKey,
		services.NewDefaultDeviceService(
			deviceRepo,
			locationRepo,
		),
		services.NewDefaultLocationService(
			locationRepo,
			params.LocationHistoryLimit,
		),
	)

	t.Cleanup(func() {
		err := client.Disconnect(ctx)
		require.NoError(t, err, "Failed to close mongodb container")

		err = container.Terminate(ctx)
		require.NoError(t, err, "Failed to close mongodb client")
	})

	return router
}
