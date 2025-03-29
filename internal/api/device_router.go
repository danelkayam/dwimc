package api

import (
	api_model "dwimc/internal/api/model"
	api_utils "dwimc/internal/api/utils"
	"dwimc/internal/model"
	"dwimc/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Devices API
// GET     /api/devices/ - get user's devices
// GET     /api/devices/:device_id - get device
// POST    /api/devices/ - upsert device
// DELETE  /api/devices/:device_id - delete device

type DeviceRouter struct {
	service services.DeviceService
}

func NewDeviceRouter(service services.DeviceService) *DeviceRouter {
	return &DeviceRouter{service: service}
}

func (r *DeviceRouter) GetAll(c *gin.Context) {
	devices, err := r.service.GetAll()
	if api_utils.HandleErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, api_model.Response[[]model.Device]{
		Data:  devices,
		Error: nil,
	})
}

func (r *DeviceRouter) Get(c *gin.Context) {
	deviceID := c.Param("device_id")

	device, err := r.service.Get(deviceID)
	if api_utils.HandleErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, api_model.Response[*model.Device]{
		Data:  device,
		Error: nil,
	})
}

func (r *DeviceRouter) Create(c *gin.Context) {
	var createParams api_model.CreateDevice

	if api_utils.BindJsonOrErrorResponse(c, &createParams) {
		return
	}

	device, err := r.service.Create(createParams.Serial, createParams.Name)
	if api_utils.HandleErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, api_model.Response[*model.Device]{
		Data:  device,
		Error: nil,
	})
}

func (r *DeviceRouter) Delete(c *gin.Context) {
	deviceID := c.Param("device_id")

	ok, err := r.service.Delete(deviceID)
	if api_utils.HandleErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, api_model.Response[api_model.Operation]{
		Data:  api_model.Operation{Success: ok},
		Error: nil,
	})
}
