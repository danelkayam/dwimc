package repositories

import "dwimc/internal/model"

type LocationRepository interface {
	GetLast(deviceId model.ID) *model.Location
	GetAllBy(deviceId model.ID) []model.Location
	Create(location *model.Location) (*model.Location, error)
	Delete(location *model.Location) bool
}
