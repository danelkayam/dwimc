package api

import (
	"dwimc/internal/services"

	"github.com/gin-gonic/gin"
)

func InitializeRouters(debugMode bool,
	deviceService services.DeviceService,
	locationService services.LocationService) *gin.Engine {

	deviceRouter := NewDeviceRouter(deviceService)
	locationRouter := NewLocationRouter(locationService)

	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	apiGroup := router.Group("/api")

	// setup device routes
	deviceGroup := apiGroup.Group("/device")
	deviceGroup.GET("/", deviceRouter.GetDevices)
	deviceGroup.GET("/:id", deviceRouter.GetDevice)
	deviceGroup.POST("/", deviceRouter.CreateDevice)
	deviceGroup.DELETE("/:id", deviceRouter.DeleteDevice)

	// setup location routes
	locationGroup := apiGroup.Group("/devices/:device_id/locations")
	locationGroup.GET("/", locationRouter.GetLocations)
	locationGroup.GET("/latest", locationRouter.GetLatestLocation)
	locationGroup.POST("/", locationRouter.CreateLocation)
	locationGroup.DELETE("/", locationRouter.DeleteLocations)
	locationGroup.DELETE("/:id", locationRouter.DeleteLocation)

	return router
}
