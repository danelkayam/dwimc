package services

import (
	"dwimc/internal/model"
	"dwimc/internal/repositories"
	"errors"
)

type LocationService interface {
	GetAllByDevice(deviceID string) ([]model.Location, error)
	GetLatestByDevice(deviceID string) (*model.Location, error)
	Create(deviceID string, latitude float64, longitude float64) (*model.Location, error)
	DeleteAllByDevice(deviceID string) (bool, error)
	Delete(deviceID string, id string) (bool, error)
}

type DefaultLocationService struct {
	repo repositories.LocationRepository
}

func NewDefaultLocationService(repo repositories.LocationRepository) LocationService {
	return &DefaultLocationService{repo: repo}
}

func (s *DefaultLocationService) GetAllByDevice(deviceID string) ([]model.Location, error) {
	// TODO - validate fields?
	return s.repo.GetAllByDevice(deviceID)
}

func (s *DefaultLocationService) GetLatestByDevice(deviceID string) (*model.Location, error) {
	// TODO - validate fields?
	location, err := s.repo.GetLatestByDevice(deviceID)
	// since we are not requesting for a specific location,
	// we can ignore the error, returning nothing
	if err != nil && !errors.Is(err, model.ErrItemNotFound) {
		return nil, err
	}

	return location, nil
}

func (s *DefaultLocationService) Create(deviceID string, latitude float64, longitude float64) (*model.Location, error) {
	// TODO - validate fields?
	return s.repo.Create(deviceID, latitude, longitude)
}

func (s *DefaultLocationService) DeleteAllByDevice(deviceID string) (bool, error) {
	// TODO - validate fields?
	return s.repo.DeleteAllByDevice(deviceID)
}

func (s *DefaultLocationService) Delete(deviceID string, id string) (bool, error) {
	// TODO - validate fields?
	return s.repo.Delete(deviceID, id)
}
