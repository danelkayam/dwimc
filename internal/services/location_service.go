package services

import "dwimc/internal/repositories"

type LocationService interface {
	// TODO - implement this
}

type DefaultLocationService struct {
	repo repositories.LocationRepository
}

func NewDefaultLocationService(repo repositories.LocationRepository) LocationService {
	return &DefaultLocationService{repo: repo}
}
