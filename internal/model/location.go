package model

type Location struct {
	Model
	DeviceID  ID      `db:"device_id"`
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
}

func WithLatitude(latitude float64) Field {
	return WithField("latitude", latitude)
}

func WithLongitude(longitude float64) Field {
	return WithField("longitude", longitude)
}
