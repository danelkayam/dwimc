package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Location struct {
	ID        bson.ObjectID `json:"id" bson:"_id"`
	CreatedAt time.Time     `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updatedAt"`
	DeviceID  bson.ObjectID `json:"device_id" bson:"deviceId"`
	Latitude  float64       `json:"latitude" binding:"required,latitude" bson:"latitude"`
	Longitude float64       `json:"longitude" binding:"required,longitude" bson:"longitude"`
}
