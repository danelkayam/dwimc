package services

import (
	"dwimc/internal/model"
	"dwimc/internal/repositories"
)

type DeviceService interface {
	GetDevices() ([]model.Device, error)
	GetDevice(serial string) (*model.Device, error)
	CreateDevice(serial, name string) (*model.Device, error)
	DeleteDevice(serial string) error
}

type DefaultDeviceService struct {
	repo repositories.DeviceRepository
}

func NewDefaultDeviceService(repo repositories.DeviceRepository) DeviceService {
	return &DefaultDeviceService{repo: repo}
}

func (s *DefaultDeviceService) GetDevices() ([]model.Device, error) {
	// TODO - implement this
	return nil, nil
}

func (s *DefaultDeviceService) GetDevice(serial string) (*model.Device, error) {
	// TODO - implement this
	return nil, nil
}

func (s *DefaultDeviceService) CreateDevice(serial, name string) (*model.Device, error) {
	// TODO - implement this
	return nil, nil
}

func (s *DefaultDeviceService) DeleteDevice(serial string) error {
	// TODO - implement this
	return nil
}
