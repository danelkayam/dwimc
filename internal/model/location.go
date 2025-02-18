package model

type Location struct {
	Model
	DeviceID  ID
	Latitude  float64
	Longitude float64
}
