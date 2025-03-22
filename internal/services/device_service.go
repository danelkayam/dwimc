package services

import "dwimc/internal/repositories"

type DeviceService interface {
	// TODO - implement this
}

type DefaultDeviceService struct {
	repo repositories.DeviceRepository
}

func NewDefaultDeviceService(repo repositories.DeviceRepository) DeviceService {
	return &DefaultDeviceService{repo: repo}
}
