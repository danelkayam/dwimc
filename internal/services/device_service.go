package services

import (
	"dwimc/internal/model"
	"dwimc/internal/repositories"
	"strings"

	"github.com/rs/zerolog/log"
)

type DeviceService interface {
	GetAll() ([]model.Device, error)
	Get(id string) (*model.Device, error)
	Exists(id string) (bool, error)
	Create(id, name string) (*model.Device, error)
	Delete(id string) (bool, error)
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

func (s *DefaultDeviceService) GetAll() ([]model.Device, error) {
	return s.repo.GetAll()
}

func (s *DefaultDeviceService) Get(id string) (*model.Device, error) {
	// TODO - validate fields?
	return s.repo.Get(id)
}

func (s *DefaultDeviceService) Exists(id string) (bool, error) {
	return s.repo.Exists(id)
}

func (s *DefaultDeviceService) Create(serial string, name string) (*model.Device, error) {
	// TODO - validate fields?
	return s.repo.Create(
		strings.TrimSpace(serial),
		strings.TrimSpace(name),
	)
}

func (s *DefaultDeviceService) Delete(id string) (bool, error) {
	// TODO - validate fields?
	defer func() {
		// deletes all locations associated with the device
		_, err := s.locationRepo.DeleteAllByDevice(id)
		if err != nil {
			log.Warn().
				Err(err).
				Str("id", id).
				Msgf("failed to delete locations associated with device: %s", id)
		}
	}()

	return s.repo.Delete(id)
}
