package api

import (
	api_model "dwimc/internal/api/model"
	"dwimc/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Locations API
// GET     /api/devices/:device_id/locations - get all locations
// GET     /api/devices/:device_id/locations/latest - get last known location
// POST    /api/devices/:device_id/locations - creates new location reporting (there will be limitation for last X locations)
// DELETE  /api/devices/:device_id/locations - delete all locations
// DELETE  /api/devices/:device_id/locations/:id - delete specific location

type LocationRouter struct {
	service services.LocationService
}

func NewLocationRouter(service services.LocationService) *LocationRouter {
	return &LocationRouter{service: service}
}

func (r *LocationRouter) GetAll(c *gin.Context) {
	deviceID := c.Param("device_id")

	locations, err := r.service.GetAllByDevice(deviceID)
	if handleErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, api_model.Response{
		Data:  locations,
		Error: nil,
	})
}

func (r *LocationRouter) GetLatest(c *gin.Context) {
	deviceID := c.Param("device_id")

	location, err := r.service.GetLatestByDevice(deviceID)
	if handleErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, api_model.Response{
		Data:  location,
		Error: nil,
	})
}

func (r *LocationRouter) Create(c *gin.Context) {
	deviceID := c.Param("device_id")

	var location api_model.CreateLocation

	if bindJsonOrErrorResponse(c, &location) {
		return
	}

	_, err := r.service.Create(deviceID, location.Latitude, location.Longitude)
	if handleErrorResponse(c, err) {
		return
	}

	// no need to retrieve location as ack
	c.JSON(http.StatusOK, api_model.Response{
		Data: map[string]any{
			"success": true,
		},
		Error: nil,
	})
}

func (r *LocationRouter) DeleteAll(c *gin.Context) {
	deviceID := c.Param("device_id")

	ok, err := r.service.DeleteAllByDevice(deviceID)
	if handleErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, api_model.Response{
		Data: map[string]any{
			"success": ok,
		},
		Error: nil,
	})
}

func (r *LocationRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	ok, err := r.service.Delete(id)
	if handleErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, api_model.Response{
		Data: map[string]any{
			"success": ok,
		},
		Error: nil,
	})
}
