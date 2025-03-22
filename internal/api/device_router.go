package api

import (
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

func (r *DeviceRouter) GetDevices(c *gin.Context) {
	// TODO - implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *DeviceRouter) GetDevice(c *gin.Context) {
	// TODO - implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *DeviceRouter) CreateDevice(c *gin.Context) {
	// TODO - implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *DeviceRouter) DeleteDevice(c *gin.Context) {
	// TODO - implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}
