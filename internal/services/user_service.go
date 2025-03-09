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
	validator := utils.NewFieldsValidator().
		WithField(model.WithEmail(email)).
		WithField(model.WithPassword(password)).
		WithValidator("email", validateEmail).
		WithValidator("password", validatePassword)

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
	if len(fields) == 0 {
		return nil, utils.AsError(model.ErrInvalidArgs, "Missing Fields")
	}

	validator := utils.NewFieldsValidator().
		WithFields(fields).
		WithValidator("email", validateEmail).
		WithValidator("password", validatePassword)

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

func validateEmail(value any) error {
	email, ok := value.(string)
	if !ok || email == "" {
		return utils.AsError(model.ErrInvalidArgs, "Invalid Email")
	}

	if !utils.IsValidEmail(email) {
		return utils.AsError(model.ErrInvalidArgs, "Invalid Email")
	}

	return nil
}

func validatePassword(value any) error {
	password, ok := value.(string)
	if !ok || password == "" {
		return utils.AsError(model.ErrInvalidArgs, "Invalid Password")
	}

	if !utils.IsValidPassword(password) {
		return utils.AsError(model.ErrInvalidArgs, "Invalid Password")
	}

	return nil
}
