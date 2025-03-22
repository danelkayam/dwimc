package model

import (
	"time"
)

type Location struct {
	ID        string    `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	DeviceID  string    `json:"deviceId" bson:"deviceId"`
	Latitude  float64   `json:"latitude" binding:"required,latitude" bson:"latitude"`
	Longitude float64   `json:"longitude" binding:"required,longitude" bson:"longitude"`
}
