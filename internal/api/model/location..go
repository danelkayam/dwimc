package api_model

type CreateLocation struct {
	Latitude  float64 `json:"latitude" binding:"required,latitude"`
	Longitude float64 `json:"longitude" binding:"required,longitude"`
}