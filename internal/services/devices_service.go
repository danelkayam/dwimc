package services

import (
	"dwimc/internal/model"
	"dwimc/internal/repositories"
	"dwimc/internal/utils"
)

type DeviceService interface {
	GetByID(id model.ID) (*model.Device, error)
	GetBySerial(userID model.ID, serial string) (*model.Device, error)
	GetAllByUserID(userID model.ID) ([]model.Device, error)
	Create(userId model.ID, serial string, name string, token string) (*model.Device, error)
	Update(id model.ID, fields ...model.Field) (*model.Device, error)
	Delete(id model.ID) error
	DeleteAllByUserID(userID model.ID) (int64, error)

	// TODO - revoke device token
	// TODO - revoke devices tokens
}

type DefaultDeviceService struct {
	repo repositories.DeviceRepository
}

func NewDefaultDeviceService(repo repositories.DeviceRepository) DeviceService {
	return &DefaultDeviceService{repo: repo}
}

func (s *DefaultDeviceService) GetByID(id model.ID) (*model.Device, error) {
	return s.repo.GetByID(id)
}

func (s *DefaultDeviceService) GetBySerial(userID model.ID, serial string) (*model.Device, error) {
	return s.repo.GetBySerial(userID, serial)
}

func (s *DefaultDeviceService) GetAllByUserID(userID model.ID) ([]model.Device, error) {
	return s.repo.GetAllByUserID(userID)
}

func (s *DefaultDeviceService) Create(userId model.ID, serial string, name string, token string) (*model.Device, error) {
	// TODO - implement this
	return nil, nil
}

func (s *DefaultDeviceService) Update(id model.ID, fields ...model.Field) (*model.Device, error) {
	if len(fields) == 0 {
		return nil, utils.AsError(model.ErrInvalidArgs, "Missing Fields")
	}

	return nil, nil
}

func (s *DefaultDeviceService) Delete(id model.ID) error {
	if err := s.repo.Delete(id); err != nil {
		return utils.AsError(err, "Failed to delete device")
	}

	return nil
}

func (s *DefaultDeviceService) DeleteAllByUserID(userID model.ID) (int64, error) {
	deleted, err := s.repo.DeleteAllByUserID(userID)
	if err != nil {
		return 0, utils.AsError(err, "Failed to delete device")
	}

	return deleted, nil
}
