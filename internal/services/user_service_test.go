package services_test

import (
	"dwimc/internal/model"
	"dwimc/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService(t *testing.T) {
	t.Run("create user", func(t *testing.T) {
		t.Run("valid user", func(t *testing.T) {
			t.Parallel()

			mockRepo := new(MockUserRepository)
			service := services.NewDefaultUserService(mockRepo)

			email := "moshe@dwimc.awesome"
			password := "DwimcPassword1!"

			mockUser := &model.User{
				Model: model.Model{
					ID:        model.ID(1),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Email:    email,
				Password: password,
			}

			mockRepo.On("Create", email, mock.AnythingOfType("string")).Return(mockUser, nil)

			user, err := service.Create(email, password)
			assert.NoErrorf(t, err, "Create User failed: %v", err)
			assert.NotNilf(t, user, "Create User failed - user is nil")
			assert.Equalf(t, user.Email, email, "Email mismatch")
			mockRepo.AssertExpectations(t)
		})

		t.Run("invalid user", func(t *testing.T) {
			email := "moshe@dwimc.awesome"
			password := "DwimcPassword1!"

			mockUser := &model.User{
				Model: model.Model{
					ID:        model.ID(1),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				Email:    email,
				Password: password,
			}

			t.Run("empty fields", func(t *testing.T) {
				t.Parallel()

				mockRepo := new(MockUserRepository)
				service := services.NewDefaultUserService(mockRepo)

				mockRepo.On("Create", email, password, mock.AnythingOfType("string")).Return(mockUser, nil)

				user, err := service.Create("", "")
				assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected error for missing fields")
				assert.Nilf(t, user, "Expected nil user, got: %v", user)
				mockRepo.AssertNotCalled(t, "Create")
			})

			t.Run("empty email", func(t *testing.T) {
				t.Parallel()

				mockRepo := new(MockUserRepository)
				service := services.NewDefaultUserService(mockRepo)

				mockRepo.On("Create", email, password, mock.AnythingOfType("string")).Return(mockUser, nil)

				user, err := service.Create("", password)
				assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected error for missing fields")
				assert.Nilf(t, user, "Expected nil user, got: %v", user)
				mockRepo.AssertNotCalled(t, "Create")
			})

			t.Run("empty password", func(t *testing.T) {
				t.Parallel()

				mockRepo := new(MockUserRepository)
				service := services.NewDefaultUserService(mockRepo)

				mockRepo.On("Create", email, password, mock.AnythingOfType("string")).Return(mockUser, nil)

				user, err := service.Create(email, "")
				assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected error for missing fields")
				assert.Nilf(t, user, "Expected nil user, got: %v", user)
				mockRepo.AssertNotCalled(t, "Create")
			})
		})
	})

	t.Run("update user", func(t *testing.T) {

		email := "moshe@dwimc.awesome"
		password := "DwimcPassword1!"

		mockUser := &model.User{
			Model: model.Model{
				ID:        model.ID(1),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Email:    email,
			Password: password,
		}

		t.Run("valid user", func(t *testing.T) {
			t.Parallel()

			mockRepo := new(MockUserRepository)
			service := services.NewDefaultUserService(mockRepo)

			mockRepo.On("Update", model.ID(1), mock.Anything).Return(mockUser, nil)

			user, err := service.Update(mockUser.ID,
				model.WithEmail(email),
				model.WithPassword(password),
			)
			assert.NoErrorf(t, err, "Update User failed: %v", err)
			assert.NotNilf(t, user, "Update User failed - user is nil")
			assert.Equalf(t, user.Email, email, "Email mismatch")
			mockRepo.AssertExpectations(t)
		})

		t.Run("invalid user", func(t *testing.T) {
			t.Run("no fields", func(t *testing.T) {
				t.Parallel()

				mockRepo := new(MockUserRepository)
				service := services.NewDefaultUserService(mockRepo)

				mockRepo.On("Update", model.ID(1), mock.Anything).Return(mockUser, nil)

				user, err := service.Update(mockUser.ID)
				assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected error for missing fields")
				assert.Nilf(t, user, "Expected nil user, got: %v", user)
				mockRepo.AssertNotCalled(t, "Update")
			})

			t.Run("invalid email", func(t *testing.T) {
				t.Parallel()

				mockRepo := new(MockUserRepository)
				service := services.NewDefaultUserService(mockRepo)

				mockRepo.On("Update", model.ID(1), mock.Anything).Return(mockUser, nil)

				user, err := service.Update(mockUser.ID, model.WithEmail("asldkfmasldkfm"))
				assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected error for invalid arguments")
				assert.Nilf(t, user, "Expected nil user, got: %v", user)
				mockRepo.AssertNotCalled(t, "Update")
			})

			t.Run("invalid password", func(t *testing.T) {
				t.Parallel()

				mockRepo := new(MockUserRepository)
				service := services.NewDefaultUserService(mockRepo)

				mockRepo.On("Update", model.ID(1), mock.Anything).Return(mockUser, nil)

				user, err := service.Update(mockUser.ID, model.WithPassword("asldkfmasldkfm"))
				assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected error for invalid arguments")
				assert.Nilf(t, user, "Expected nil user, got: %v", user)
				mockRepo.AssertNotCalled(t, "Update")
			})
		})
	})

	t.Run("delete user", func(t *testing.T) {
		t.Parallel()

		mockRepo := new(MockUserRepository)
		service := services.NewDefaultUserService(mockRepo)

		mockRepo.On("Delete", model.ID(1)).Return(nil)

		err := service.Delete(1)
		assert.NoErrorf(t, err, "Update User failed: %v", err)
		mockRepo.AssertExpectations(t)
	})
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(id model.ID) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Create(email, hashedPassword string) (*model.User, error) {
	args := m.Called(email, hashedPassword)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(id model.ID, fields ...model.Field) (*model.User, error) {
	args := m.Called(id, fields)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id model.ID) error {
	args := m.Called(id)
	return args.Error(0)
}
