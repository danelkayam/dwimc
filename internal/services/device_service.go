package services

import (
	"dwimc/internal/model"
	"dwimc/internal/repositories"

	"github.com/rs/zerolog/log"
)

type DeviceService interface {
	GetDevices() ([]model.Device, error)
	GetDevice(serial string) (*model.Device, error)
	CreateDevice(serial, name string) (*model.Device, error)
	DeleteDevice(serial string) (bool, error)
}

type DefaultDeviceService struct {
	repo         repositories.DeviceRepository
	locationRepo repositories.LocationRepository
}

func NewDefaultDeviceService(
	repo repositories.DeviceRepository,
	locationRepo repositories.LocationRepository,
) DeviceService {
	return &DefaultDeviceService{
		repo:         repo,
		locationRepo: locationRepo,
	}
}

func (s *DefaultDeviceService) GetDevices() ([]model.Device, error) {
	return s.repo.GetDevices()
}

func (s *DefaultDeviceService) GetDevice(serial string) (*model.Device, error) {
	// TODO - fields validation?
	return s.repo.GetDevice(serial)
}

func (s *DefaultDeviceService) CreateDevice(serial, name string) (*model.Device, error) {
	// TODO - fields validation?
	return s.repo.CreateDevice(serial, name)
}

func (s *DefaultDeviceService) DeleteDevice(serial string) (bool, error) {
	defer func() {
		// deletes all locations associated with the device
		_, err := s.locationRepo.DeleteLocations(serial)
		if err != nil {
			log.Warn().
				Err(err).
				Str("serial", serial).
				Msg("failed to delete locations associated with the device")

		}
	}()

	return s.repo.DeleteDevice(serial)
}
