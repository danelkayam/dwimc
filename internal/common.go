package service

import "time"

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

type Operation struct {
	Success bool `json:"success"`
}

type UpdateParams struct {
	Serial   string   `json:"serial" binding:"required,min=4,max=64"`
	Name     string   `json:"name" binding:"required,max=64"`
	Location Location `json:"location" binding:"required"`
}
