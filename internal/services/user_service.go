package services

import (
	"dwimc/internal/model"
	"dwimc/internal/repositories"
	"dwimc/internal/utils"
)

type UserService interface {
	Create(email string, password string) (*model.User, error)
	Update(id model.ID, fields ...model.Field) (*model.User, error)
	Delete(id model.ID) error
}

type DefaultUserService struct {
	repo repositories.UserRepository
}

func NewDefaultUserService(repo repositories.UserRepository) UserService {
	return &DefaultUserService{repo: repo}
}

func (s *DefaultUserService) Create(email string, password string) (*model.User, error) {
	validator := utils.NewWithValidator().
		WithField(model.WithEmail(email)).
		WithField(model.WithPassword(password)).
		WithValidator(emailValidator()).
		WithValidator(passwordValidator())

	if err := validator.Validate(); err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, utils.AsError(model.ErrInternal, err.Error())
	}

	user, err := s.repo.Create(email, hashedPassword)
	if err != nil {
		return nil, utils.AsError(err, "Failed to create user")
	}

	return user, nil
}

func (s *DefaultUserService) Update(id model.ID, fields ...model.Field) (*model.User, error) {
	validator := utils.NewWithValidator().
		WithFields(fields).
		WithValidator(emailValidator()).
		WithValidator(passwordValidator()).
		WithNoFieldsValidation(utils.AsError(model.ErrInvalidArgs, "Missing Fields")).
		WithStrictModeValidation(utils.AsError(model.ErrInvalidArgs, "Invalid Fields"))

	if err := validator.Validate(); err != nil {
		return nil, err
	}

	updated, err := s.repo.Update(id, fields...)
	if err != nil {
		return nil, utils.AsError(err, "Failed to update user")
	}

	return updated, nil
}

func (s *DefaultUserService) Delete(id model.ID) error {
	if err := s.repo.Delete(id); err != nil {
		return utils.AsError(err, "Failed to delete user")
	}

	return nil
}

func emailValidator() utils.Validator {
	return utils.WithFieldValidator(
		"email",
		"required,email,min=5,max=254",
		"Invalid Email",
	)
}

func passwordValidator() utils.Validator {
	return utils.WithFieldValidator(
		"password",
		"required,min=8,max=64,strong_password",
		"Invalid Password",
	)
}
