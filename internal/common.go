package service

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data  any            `json:"data"`
	Error *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type Operation struct {
	Success bool `json:"success"`
}

type Location struct {
	Latitude  float64 `json:"latitude" binding:"required,latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" binding:"required,longitude" bson:"longitude"`
}

type Device struct {
	Serial    string    `json:"serial" bson:"serial"`
	Name      string    `json:"name" bson:"name"`
	Location  Location  `json:"location" bson:"location"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type UpdateParams struct {
	Serial string `json:"serial" binding:"required,min=4,max=64"`
	// todo - Name should not be required!
	Name     string   `json:"name" binding:"required,max=64"`
	Location Location `json:"location" binding:"required"`
}

// Helper functions

// handleInternalError handles Internal Server Error response if the given err argument is not nil.
// returns true if an error response was sent back and calling function should be terminate,
// false otherwise.
func handleInternalError(err error, message string, c *gin.Context) bool {
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
