package services

import (
	"dwimc/internal/model"
	"dwimc/internal/repositories"
	"dwimc/internal/utils"
)

type UserService interface {
	Create(email string, password string) (*model.User, error)
	Update(id model.ID, fields ...model.UpdateField) (*model.User, error)
	Delete(id model.ID) error
}

type DefaultUserService struct {
	repo repositories.UserRepository
}

func NewDefaultUserService(repo repositories.UserRepository) UserService {
	return &DefaultUserService{repo: repo}
}

func (s *DefaultUserService) Create(email string, password string) (*model.User, error) {
	if email == "" {
		return nil, utils.AsError(model.ErrInvalidArgs, "Missing Email")
	}

	if !utils.IsValidEmail(email) {
		return nil, utils.AsError(model.ErrInvalidArgs, "Invalid Email")
	}

	if password == "" {
		return nil, utils.AsError(model.ErrInvalidArgs, "Missing Password")
	}

	if !utils.IsValidPassword(password) {
		return nil, utils.AsError(model.ErrInvalidArgs, "Invalid Password")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, utils.AsError(model.ErrInternal, err.Error())
	}

	user, err := s.repo.Create(email, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *DefaultUserService) Update(id model.ID, fields ...model.UpdateField) (*model.User, error) {
	if len(fields) == 0 {
		return nil, utils.AsError(model.ErrInvalidArgs, "Missing Fields")
	}

	if err := utils.NewUpdateFieldsValidator(fields).
		WithValidator("email", validateEmail).
		WithValidator("password", validatePassword).
		Validate(); err != nil {

		return nil, err
	}

	updated, err := s.repo.Update(id, fields...)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *DefaultUserService) Delete(id model.ID) error {
	return s.repo.Delete(id)
}

func validateEmail(value any) error {
	email, ok := value.(string)
	if !ok {
		return utils.AsError(model.ErrInvalidArgs, "Invalid Email")
	}

	if !utils.IsValidEmail(email) {
		return utils.AsError(model.ErrInvalidArgs, "Invalid Email")
	}

	return nil
}

func validatePassword(value any) error {
	password, ok := value.(string)
	if !ok {
		return utils.AsError(model.ErrInvalidArgs, "Invalid Password")
	}

	if !utils.IsValidPassword(password) {
		return utils.AsError(model.ErrInvalidArgs, "Invalid Password")
	}

	return nil
}
