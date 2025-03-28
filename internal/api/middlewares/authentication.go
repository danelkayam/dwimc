package middlewares

import (
	api_utils "dwimc/internal/api/utils"
	"dwimc/internal/model"

	"github.com/gin-gonic/gin"
)

func ApiKeyAuthenticationMiddleware(secretAPIKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get("X-API-Key")

		if apiKey != secretAPIKey {
			if api_utils.HandleErrorResponse(c, model.ErrUnauthenticated) {
				return
			}
		}

		c.Next()
	}
}