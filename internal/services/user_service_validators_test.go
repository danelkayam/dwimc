package services

import (
	"dwimc/internal/model"
	"dwimc/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserServiceValidators(t *testing.T) {
	t.Run("validate email", func(t *testing.T) {
		t.Run("valid one", func(t *testing.T) {
			t.Parallel()

			validator := utils.NewWithValidator().
				WithField(model.WithEmail("mekmek@dwimc.awesome")).
				WithValidator(emailValidator())

			err := validator.Validate()
			assert.NoErrorf(t, err, "Expected valid email: %v", err)
		})

		t.Run("invalid", func(t *testing.T) {
			t.Parallel()

			invalidEmails := []string{
				"mekmek@@@dwimc.awesome",
				"mekmek",
				"dwimc.awesome",
				"a@b",
				"m",
				"",
				"1111",
				"nil",
				"verylonglocalpart_exceeding_sixty_four_characters_abcdefghijklmnopqrstuvwxyz0123456789_abcdefghijklmnopqrstuvwxyz@extremely-long-and-fake-domain-to-test-boundary-conditions-and-validate-the-limits-of-the-email-address-fields-in-your-application-to-ensure-security-against-attacks.com",
			}

			for _, invalid := range invalidEmails {
				validator := utils.NewWithValidator().
					WithField(model.WithEmail(invalid)).
					WithValidator(emailValidator())

				err := validator.Validate()
				assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected invalid email")
			}
		})
	})

	t.Run("validate password", func(t *testing.T) {
		t.Run("valid one", func(t *testing.T) {
			t.Parallel()

			validator := utils.NewWithValidator().
				WithField(model.WithPassword("Dw1mcAw3some&&")).
				WithValidator(passwordValidator())

			err := validator.Validate()
			assert.NoErrorf(t, err, "Expected valid password: %v", err)
		})

		t.Run("invalid", func(t *testing.T) {
			t.Parallel()

			invalidPassword := []string{
				"DwimcAwesome",
				"Mek123456",
				"Mek!@#$%^",
				"Mek1!",
				"123",
				"",
				"nil",
				"M9#dT7zW$k@1hXqE^6b!2pV3fU%4yG&sR*0jL8cQ+5nI?oA-CwDrZtHxuFgYKmPe!",
			}

			for _, invalid := range invalidPassword {
				validator := utils.NewWithValidator().
					WithField(model.WithPassword(invalid)).
					WithValidator(passwordValidator())

				err := validator.Validate()
				assert.ErrorIsf(t, err, model.ErrInvalidArgs, "Expected invalid password")
			}
		})
	})
}
