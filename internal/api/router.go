package api

import (
	"dwimc/internal/services"

	"github.com/gin-gonic/gin"
)

func InitializeRouters(
	debugMode bool,
	userService services.UserService,
	deviceService services.DeviceService,
	locationService services.LocationService) *gin.Engine {

	userRouter := NewUserRouter(userService)
	deviceRouter := NewDeviceRouter(deviceService)
	locationRouter := NewLocationRouter(locationService)

	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	apiGroup := router.Group("/api")
	// setup user routes
	userGroup := apiGroup.Group("/users")
	userGroup.POST("/auth/register", userRouter.RegisterUser)
	userGroup.POST("/auth/signin", userRouter.SignInUser)
	userGroup.GET("/me", userRouter.GetCurrentUser)
	userGroup.PUT("/", userRouter.UpdateCurrentUser)
	userGroup.DELETE("/", userRouter.DeleteCurrentUser)

	// setup device routes
	deviceGroup := apiGroup.Group("/devices")
	deviceGroup.GET("/", deviceRouter.GetDevices)
	deviceGroup.POST("/", deviceRouter.CreateDevice)
	deviceGroup.GET("/:device_id", deviceRouter.GetDevice)
	deviceGroup.PUT("/:device_id", deviceRouter.UpdateDevice)
	deviceGroup.DELETE("/:device_id", deviceRouter.DeleteDevice)
	deviceGroup.DELETE("/:device_id/token", deviceRouter.RevokeDeviceToken)
	deviceGroup.DELETE("/tokens", deviceRouter.RevokeAllDeviceTokens)

	// setup location routes
	locationGroup := apiGroup.Group("/devices/:device_id/locations")
	locationGroup.GET("/", locationRouter.GetAllLocations)
	locationGroup.GET("/latest", locationRouter.GetLatestLocation)
	locationGroup.POST("/", locationRouter.CreateLocation)
	locationGroup.DELETE("/", locationRouter.DeleteAllLocations)
	locationGroup.DELETE("/:location_id", locationRouter.DeleteLocation)

	return router
}
