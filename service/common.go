package service

import "time"

type Position struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Device struct {
	Serial    string    `json:"serial"`
	Name      string    `json:"name"`
	Position  Position  `json:"position"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateParams struct {
	Serial   string   `json:"serial"`
	Name     string   `json:"name"`
	Position Position `json:"position"`
}
