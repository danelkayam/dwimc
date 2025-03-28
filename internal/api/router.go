package api

import (
	"dwimc/internal/api/middlewares"
	"dwimc/internal/services"
	"dwimc/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitializeRouters(debugMode bool,
	deviceService services.DeviceService,
	locationService services.LocationService,
) *gin.Engine {

	deviceRouter := NewDeviceRouter(deviceService)
	locationRouter := NewLocationRouter(locationService)

	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		utils.RegisterValidations(v)
	}

	router := gin.Default()
	apiGroup := router.Group("/api")

	// setup device routes
	deviceGroup := apiGroup.Group("/devices")
	deviceGroup.GET("/", deviceRouter.GetAll)
	deviceGroup.GET("/:device_id", deviceRouter.Get)
	deviceGroup.POST("/", deviceRouter.Create)
	deviceGroup.DELETE("/:device_id", deviceRouter.Delete)

	// setup location routes
	locationGroup := deviceGroup.Group("/:device_id/locations")
	locationGroup.Use(middlewares.DeviceExistsMiddleware(deviceService))

	locationGroup.GET("/", locationRouter.GetAll)
	locationGroup.GET("/latest", locationRouter.GetLatest)
	locationGroup.POST("/", locationRouter.Create)
	locationGroup.DELETE("/", locationRouter.DeleteAll)
	locationGroup.DELETE("/:id", locationRouter.Delete)

	return router
}
