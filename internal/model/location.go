package model

type Location struct {
	Model
	DeviceID  ID      `db:"device_id"`
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
}
