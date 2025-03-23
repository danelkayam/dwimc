package services

import (
	"dwimc/internal/model"
	"dwimc/internal/repositories"
	"errors"
)

type LocationService interface {
	GetLocations(serial string) ([]model.Location, error)
	GetLatestLocation(serial string) (*model.Location, error)
	CreateLocation(location model.Location) (*model.Location, error)
	DeleteLocation(id string) (bool, error)
}

type DefaultLocationService struct {
	repo repositories.LocationRepository
}

func NewDefaultLocationService(repo repositories.LocationRepository) LocationService {
	return &DefaultLocationService{repo: repo}
}

func (s *DefaultLocationService) GetLocations(serial string) ([]model.Location, error) {
	return s.repo.GetLocations(serial)
}

func (s *DefaultLocationService) GetLatestLocation(serial string) (*model.Location, error) {
	location, err := s.repo.GetLatestLocation(serial)
	// since we are not requesting for a specific location,
	// we can ignore the error, returning nothing
	if err != nil && !errors.Is(err, model.ErrItemNotFound) {
		return nil, err
	}

	return location, nil
}

func (s *DefaultLocationService) CreateLocation(location model.Location) (*model.Location, error) {
	// TODO - validate location?
	return s.repo.CreateLocation(location)
}

func (s *DefaultLocationService) DeleteLocation(id string) (bool, error) {
	return s.repo.DeleteLocation(id)
}
