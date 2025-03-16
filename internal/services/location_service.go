package services

import (
	"dwimc/internal/model"
	"dwimc/internal/repositories"
	"dwimc/internal/utils"
)

type LocationService interface {
	GetLastByDeviceID(deviceID model.ID) (*model.Location, error)
	GetAllByDeviceID(deviceID model.ID) ([]model.Location, error)
	Create(deviceID model.ID, latitude float64, longitude float64) (*model.Location, error)
	Delete(id model.ID) error
	DeleteAllBy(deviceID model.ID) (int64, error)
}

type DefaultLocationService struct {
	repo repositories.LocationRepository
}

func NewDefaultLocationService(repo repositories.LocationRepository) LocationService {
	return &DefaultLocationService{
		repo: repo,
	}
}

func (s *DefaultLocationService) GetLastByDeviceID(deviceID model.ID) (*model.Location, error) {
	return s.repo.GetLastByDeviceID(deviceID)
}

func (s *DefaultLocationService) GetAllByDeviceID(deviceID model.ID) ([]model.Location, error) {
	return s.repo.GetAllByDeviceID(deviceID)
}

func (s *DefaultLocationService) Create(deviceID model.ID, latitude float64, longitude float64) (*model.Location, error) {
	validator := utils.NewWithValidator().
		WithField(model.WithLatitude(latitude)).
		WithField(model.WithLongitude(longitude)).
		WithValidator(latitudeValidator()).
		WithValidator(longitudeValidator())

	if err := validator.Validate(); err != nil {
		return nil, err
	}

	location, err := s.repo.Create(deviceID, latitude, longitude)
	if err != nil {
		return nil, utils.AsError(err, "Failed to create location")
	}

	return location, nil
}

func (s *DefaultLocationService) Delete(id model.ID) error {
	err := s.repo.Delete(id)
	if err != nil {
		return utils.AsError(err, "Failed to delete location")
	}

	return nil
}

func (s *DefaultLocationService) DeleteAllBy(deviceID model.ID) (int64, error) {
	deleted, err := s.repo.DeleteAllBy(deviceID)
	if err != nil {
		return 0, utils.AsError(err, "Failed to delete locations")
	}

	return deleted, nil
}

func latitudeValidator() utils.Validator {
	return utils.WithFieldValidator(
		"latitude",
		"required,latitude",
		"Invalid Latitude",
	)
}

func longitudeValidator() utils.Validator {
	return utils.WithFieldValidator(
		"longitude",
		"required,longitude",
		"Invalid Longitude",
	)
}
