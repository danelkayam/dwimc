package repositories

import "dwimc/internal/model"

type DeviceRepository interface {
	Get(id model.ID) *model.Device
	GetAllBy(userID model.ID) []model.Device
	Create(device *model.Device) (*model.Device, error)
	Update(device *model.Device) bool
	Delete(device *model.Device) bool
}
