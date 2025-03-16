package services

import (
	"dwimc/internal/model"
	"dwimc/internal/repositories"
	"dwimc/internal/utils"

	"github.com/google/uuid"
)

type DeviceService interface {
	GetByID(id model.ID) (*model.Device, error)
	GetBySerial(userID model.ID, serial string) (*model.Device, error)
	GetAllByUserID(userID model.ID) ([]model.Device, error)
	Create(userId model.ID, serial string, name string) (*model.Device, error)
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

func (s *DefaultDeviceService) Create(userId model.ID, serial string, name string) (*model.Device, error) {
	validator := utils.NewWithValidator().
		WithField(model.WithSerial(serial)).
		WithField(model.WithName(name)).
		WithValidator(serialValidator()).
		WithValidator(nameValidator())

	if err := validator.Validate(); err != nil {
		return nil, err
	}

	token := uuid.NewString()

	device, err := s.repo.Create(userId, serial, name, token)
	if err != nil {
		return nil, utils.AsError(err, "Failed to create device")
	}

	return device, nil
}

func (s *DefaultDeviceService) Update(id model.ID, fields ...model.Field) (*model.Device, error) {
	validator := utils.NewWithValidator().
		WithFields(fields).
		WithValidator(serialValidator()).
		WithValidator(nameValidator()).
		WithNoFieldsValidation(utils.AsError(model.ErrInvalidArgs, "Missing Fields")).
		WithStrictModeValidation(utils.AsError(model.ErrInvalidArgs, "Invalid Fields"))

	if err := validator.Validate(); err != nil {
		return nil, err
	}

	device, err := s.repo.Update(id, fields...)
	if err != nil {
		return nil, utils.AsError(err, "Failed to update device")
	}

	return device, nil
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
		return 0, utils.AsError(err, "Failed to delete devices")
	}

	return deleted, nil
}

func serialValidator() utils.Validator {
	return utils.WithFieldValidator(
		"serial",
		"required,min=8,max=64",
		"Invalid Serial",
	)
}

func nameValidator() utils.Validator {
	return utils.WithFieldValidator(
		"name",
		"required,min=1,max=254",
		"Invalid Name",
	)
}
