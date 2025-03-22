package api

import (
	"dwimc/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Requires: user_token (JWT)
// GET     /api/devices/ - get user's devices
// POST    /api/devices/ - creates new device
// GET     /api/devices/:device_id - get device
// PUT     /api/devices/:device_id - update device
// DELETE  /api/devices/:device_id - delete device
// DELETE  /api/devices/:device_id/token - revoke specific device's token
// DELETE  /api/devices/tokens  - revokes all user's devices tokens (removes them)

type DeviceRouter struct {
	deviceService services.DeviceService
}

func NewDeviceRouter(deviceService services.DeviceService) *DeviceRouter {
	return &DeviceRouter{deviceService: deviceService}
}

func (r *DeviceRouter) GetDevices(c *gin.Context) {
	// TODO implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *DeviceRouter) CreateDevice(c *gin.Context) {
	// TODO implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *DeviceRouter) GetDevice(c *gin.Context) {
	// TODO implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *DeviceRouter) UpdateDevice(c *gin.Context) {
	// TODO implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *DeviceRouter) DeleteDevice(c *gin.Context) {
	// TODO implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *DeviceRouter) RevokeDeviceToken(c *gin.Context) {
	// TODO implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *DeviceRouter) RevokeAllDeviceTokens(c *gin.Context) {
	// TODO implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}
