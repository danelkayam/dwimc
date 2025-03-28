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
	deviceGroup.GET("/", deviceRouter.GetAll)
	deviceGroup.GET("/:id", deviceRouter.Get)
	deviceGroup.POST("/", deviceRouter.Create)
	deviceGroup.DELETE("/:id", deviceRouter.Delete)

	// setup location routes
	locationGroup := apiGroup.Group("/devices/:device_id/locations")
	locationGroup.GET("/", locationRouter.GetAll)
	locationGroup.GET("/latest", locationRouter.GetLatest)
	locationGroup.POST("/", locationRouter.Create)
	locationGroup.DELETE("/", locationRouter.DeleteAll)
	locationGroup.DELETE("/:id", locationRouter.Delete)

	return router
}
