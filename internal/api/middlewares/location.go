package middlewares

import (
	api_utils "dwimc/internal/api/utils"
	"dwimc/internal/model"
	"dwimc/internal/services"

	"github.com/gin-gonic/gin"
)

func DeviceExistsMiddleware(service services.DeviceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deviceID := c.Param("device_id")

		exists, err := service.Exists(deviceID)
		if api_utils.HandleErrorResponse(c, err) {
			return
		}

		if !exists {
			if api_utils.HandleErrorResponse(c, model.ErrItemNotFound) {
				return
			}
		}

		c.Next()
	}
}
