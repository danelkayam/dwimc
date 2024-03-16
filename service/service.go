package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service struct {
	Store        *Store
	SecretApiKey string
	server       *http.Server
}

type Response struct {
	Data  interface{}    `json:"data"`
	Error *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func (service *Service) Start(port string) error {
	router := gin.Default()

	apiRouter := router.Group("/api")
	apiRouter.Use(validateApiKey(service.SecretApiKey))

	apiRouter.GET("/devices/:device", service.handleGet)
	apiRouter.GET("/devices", service.handleGetAll)
	apiRouter.POST("/devices", service.handlePost)

	service.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	log.Printf("Lifting service on port: %v\n", port)

	return service.server.ListenAndServe()
}

func (service *Service) Stop(ctx context.Context) error {
	return service.server.Shutdown(ctx)
}

func validateApiKey(secretApiKey string) gin.HandlerFunc {
	if len(secretApiKey) > 0 {
		return func(c *gin.Context) {
			apiKey := c.Request.Header.Get("X-API-Key")

			if apiKey != secretApiKey {
				c.AbortWithStatusJSON(http.StatusUnauthorized, Response{
					Data:  nil,
					Error: &ErrorResponse{Message: "Unauthorized!"},
				})

				return
			}

			c.Next()
		}
	}

	return func(c *gin.Context) { c.Next() }
}

func (service *Service) handlePost(c *gin.Context) {
	var params UpdateParams

	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Data:  nil,
			Error: &ErrorResponse{Message: "Invalid request body!"},
		})
		return
	}

	// TODO: - validate params structure

	operation, err := service.Store.Upsert(params)

	if service.handleDatabaseError(err, "Failed upserting device", c) {
		return
	}

	c.JSON(http.StatusOK, Response{
		Data:  operation,
		Error: nil,
	})
}

func (service *Service) handleGet(c *gin.Context) {
	deviceId := c.Param("device")
	device, err := service.Store.GetOne(deviceId)

	if service.handleDatabaseError(err, "Failed getting device", c) {
		return
	}

	if device == nil {
		c.JSON(http.StatusNotFound, Response{
			Data:  nil,
			Error: &ErrorResponse{Message: "Device not found!"},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data:  device,
		Error: nil,
	})
}

func (service *Service) handleGetAll(c *gin.Context) {
	devices, err := service.Store.GetAll()

	if service.handleDatabaseError(err, "Failed getting devices", c) {
		return
	}

	c.JSON(http.StatusOK, Response{
		Data:  devices,
		Error: nil,
	})
}

// handleDatabaseError handles Internal Server Error response if the given err argument is not nil.
// returns true if an error response was sent back and calling function should be terminate,
// false otherwise.
func (service *Service) handleDatabaseError(err error, message string, c *gin.Context) bool {
	if err != nil {
		log.Printf("%s: %s\n", message, err)

		c.JSON(http.StatusInternalServerError, Response{
			Data:  nil,
			Error: &ErrorResponse{Message: "Something went wrong!"},
		})
		return true
	}

	return false
}
