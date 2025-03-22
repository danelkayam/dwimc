package services

import (
	"dwimc/internal/model"
	"dwimc/internal/repositories"
)

type LocationService interface {
	GetLocations(serial string) ([]model.Location, error)
	GetLatestLocation(serial string) (*model.Location, error)
	CreateLocation(location model.Location) (*model.Location, error)
	DeleteLocation(id string) error
}

type DefaultLocationService struct {
	repo repositories.LocationRepository
}

func NewDefaultLocationService(repo repositories.LocationRepository) LocationService {
	return &DefaultLocationService{repo: repo}
}

func (s *DefaultLocationService) GetLocations(serial string) ([]model.Location, error) {
	// TODO - implement this
	return nil, nil
}

func (s *DefaultLocationService) GetLatestLocation(serial string) (*model.Location, error) {
	// TODO - implement this
	return nil, nil
}

func (s *DefaultLocationService) CreateLocation(location model.Location) (*model.Location, error) {
	// TODO - implement this
	return nil, nil
}

func (s *DefaultLocationService) DeleteLocation(id string) error {
	// TODO - implement this
	return nil
}
