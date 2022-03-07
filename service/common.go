package service

import "time"

type Position struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type Device struct {
	Serial    string    `json:"serial" bson:"serial"`
	Name      string    `json:"name" bson:"name"`
	Position  Position  `json:"position" bson:"position"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type Operation struct {
	Success bool `json:"success"`
}

type UpdateParams struct {
	Serial   string   `json:"serial"`
	Name     string   `json:"name"`
	Position Position `json:"position"`
}
