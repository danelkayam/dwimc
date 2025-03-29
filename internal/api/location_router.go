package api

import (
	api_model "dwimc/internal/api/model"
	api_utils "dwimc/internal/api/utils"
	"dwimc/internal/model"
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
	if api_utils.HandleErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, api_model.Response[[]model.Location]{
		Data:  locations,
		Error: nil,
	})
}

func (r *LocationRouter) GetLatest(c *gin.Context) {
	deviceID := c.Param("device_id")

	location, err := r.service.GetLatestByDevice(deviceID)
	if api_utils.HandleErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, api_model.Response[*model.Location]{
		Data:  location,
		Error: nil,
	})
}

func (r *LocationRouter) Create(c *gin.Context) {
	deviceID := c.Param("device_id")

	var location api_model.CreateLocation

	if api_utils.BindJsonOrErrorResponse(c, &location) {
		return
	}

	_, err := r.service.Create(deviceID, location.Latitude, location.Longitude)
	if api_utils.HandleErrorResponse(c, err) {
		return
	}

	// no need to retrieve location as ack
	c.JSON(http.StatusOK, api_model.Response[api_model.Operation]{
		Data:  api_model.Operation{Success: true},
		Error: nil,
	})
}

func (r *LocationRouter) DeleteAll(c *gin.Context) {
	deviceID := c.Param("device_id")

	ok, err := r.service.DeleteAllByDevice(deviceID)
	if api_utils.HandleErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, api_model.Response[api_model.Operation]{
		Data:  api_model.Operation{Success: ok},
		Error: nil,
	})
}

func (r *LocationRouter) Delete(c *gin.Context) {
	deviceID := c.Param("device_id")
	id := c.Param("id")

	ok, err := r.service.Delete(deviceID, id)
	if api_utils.HandleErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, api_model.Response[api_model.Operation]{
		Data:  api_model.Operation{Success: ok},
		Error: nil,
	})
}
