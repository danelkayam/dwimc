package api

import (
	"dwimc/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Locations API
// GET     /api/devices/:device_id/locations - get all locations
// GET     /api/devices/:device_id/locations/latest - get last known location
// POST    /api/devices/:device_id/locations - creates new location reporting (there will be limitation for last X locations)
// DELETE  /api/devices/:device_id/locations - delete all locations
// DELETE  /api/devices/:device_id/locations/:location_id - delete specific location

type LocationRouter struct {
	service services.LocationService
}

func NewLocationRouter(service services.LocationService) *LocationRouter {
	return &LocationRouter{service: service}
}

func (r *LocationRouter) GetLocations(c *gin.Context) {
	// TODO - implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *LocationRouter) GetLatestLocation(c *gin.Context) {
	// TODO - implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *LocationRouter) CreateLocation(c *gin.Context) {
	// TODO - implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *LocationRouter) DeleteLocations(c *gin.Context) {
	// TODO - implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}

func (r *LocationRouter) DeleteLocation(c *gin.Context) {
	// TODO - implement this
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented!",
	})
}
