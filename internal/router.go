package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DevicesRouter struct {
	repository DevicesRepository
}

func CreateDevicesRouter(repository DevicesRepository, group *gin.RouterGroup, middlewares ...gin.HandlerFunc) DevicesRouter {
	router := &DevicesRouter{
		repository: repository,
	}
	apiRouter := group.Group("/devices")
	apiRouter.Use(middlewares...)
	apiRouter.GET("/:serial", router.handleGet)
	apiRouter.GET("/", router.handleGetAll)
	apiRouter.POST("/", router.handlePost)

	return *router
}

func (router *DevicesRouter) handleGet(c *gin.Context) {
	serial := c.Param("serial")
	device, err := router.repository.get(serial)

	if handleInternalError(err, "Failed getting device", c) {
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

func (router *DevicesRouter) handleGetAll(c *gin.Context) {
	devices, err := router.repository.getAll()

	if handleInternalError(err, "Failed getting devices", c) {
		return
	}

	c.JSON(http.StatusOK, Response{
		Data:  devices,
		Error: nil,
	})
}

func (router *DevicesRouter) handlePost(c *gin.Context) {
	var params UpdateParams

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Data:  nil,
			Error: &ErrorResponse{Message: "Invalid request body!"},
		})
		return
	}

	operation, err := router.repository.update(params)

	if handleInternalError(err, "Failed updating device", c) {
		return
	}

	c.JSON(http.StatusOK, Response{
		Data:  operation,
		Error: nil,
	})
}
