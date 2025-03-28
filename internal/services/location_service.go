package services

import (
	"dwimc/internal/model"
	"dwimc/internal/repositories"
	"errors"

	"github.com/rs/zerolog/log"
)

type LocationService interface {
	GetAllByDevice(deviceID string) ([]model.Location, error)
	GetLatestByDevice(deviceID string) (*model.Location, error)
	Create(deviceID string, latitude float64, longitude float64) (*model.Location, error)
	DeleteAllByDevice(deviceID string) (bool, error)
	Delete(deviceID string, id string) (bool, error)
}

type DefaultLocationService struct {
	repo         repositories.LocationRepository
	historyLimit int
}

func NewDefaultLocationService(
	repo repositories.LocationRepository,
	historyLimit int,
) LocationService {
	return &DefaultLocationService{
		repo:         repo,
		historyLimit: historyLimit,
	}
}

func (s *DefaultLocationService) GetAllByDevice(deviceID string) ([]model.Location, error) {
	return s.repo.GetAllByDevice(deviceID)
}

func (s *DefaultLocationService) GetLatestByDevice(deviceID string) (*model.Location, error) {
	location, err := s.repo.GetLatestByDevice(deviceID)
	// since we are not requesting for a specific location,
	// we can ignore the error, returning nothing
	if err != nil && !errors.Is(err, model.ErrItemNotFound) {
		return nil, err
	}

	return location, nil
}

func (s *DefaultLocationService) Create(deviceID string, latitude float64, longitude float64) (*model.Location, error) {
	location, err := s.repo.Create(deviceID, latitude, longitude)
	if err != nil {
		log.Warn().
			Err(err).
			Str("deviceID", deviceID).
			Float64("latitude", latitude).
			Float64("longitude", longitude).
			Msg("Failed to create location")

		return nil, err
	}

	if s.historyLimit > 0 {
		deleted, err := s.repo.DeleteOldByDevice(
			location.DeviceID.Hex(),
			s.historyLimit,
		)
		if err != nil {
			log.Warn().
				Err(err).
				Str("deviceID", location.DeviceID.Hex()).
				Int("skip", s.historyLimit).
				Msg("Failed to delete old locations")

		} else {
			log.Info().
				Str("deviceID", location.DeviceID.Hex()).
				Int("skip", s.historyLimit).
				Int64("deleted", deleted).
				Msg("Success deleting old locations")
		}
	}

	return location, nil
}

func (s *DefaultLocationService) DeleteAllByDevice(deviceID string) (bool, error) {
	return s.repo.DeleteAllByDevice(deviceID)
}

func (s *DefaultLocationService) Delete(deviceID string, id string) (bool, error) {
	return s.repo.Delete(deviceID, id)
}
