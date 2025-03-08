package utils_test

import (
	"dwimc/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidators(t *testing.T)  {
	
	t.Run("validate email", func(t *testing.T) {
		t.Run("valid one", func(t *testing.T) {
			t.Parallel()

			assert.Truef(t, utils.IsValidEmail("mekmek@dwimc.awesome"), "Expected valid email")
		})

		t.Run("invalid", func(t *testing.T) {
			t.Parallel()

			assert.Falsef(t, utils.IsValidEmail("mekmek@@@dwimc.awesome"), "Expected invalid email")
			assert.Falsef(t, utils.IsValidEmail("mekmek"), "Expected invalid email")
			assert.Falsef(t, utils.IsValidEmail("dwimc.awesome"), "Expected invalid email")
		})

		t.Run("too short", func(t *testing.T) {
			t.Parallel()

			assert.Falsef(t, utils.IsValidEmail("a@b"), "Expected invalid email")
		})

		t.Run("too long", func(t *testing.T) {
			t.Parallel()

			too_long := "verylonglocalpart_exceeding_sixty_four_characters_abcdefghijklmnopqrstuvwxyz0123456789_abcdefghijklmnopqrstuvwxyz@extremely-long-and-fake-domain-to-test-boundary-conditions-and-validate-the-limits-of-the-email-address-fields-in-your-application-to-ensure-security-against-attacks.com"
			assert.Falsef(t, utils.IsValidEmail(too_long), "Expected too long email address")
		})
	})

	t.Run("validate password", func(t *testing.T) {
		t.Run("valid one", func(t *testing.T) {
			t.Parallel()

			assert.Truef(t, utils.IsValidPassword("Dw1mcAw3some&&"), "Expected valid password")
		})

		t.Run("not strong enough", func(t *testing.T) {
			t.Parallel()

			assert.Falsef(t, utils.IsValidPassword("DwimcAwesome"), "Expected week password")
			assert.Falsef(t, utils.IsValidPassword("Mek123456"), "Expected week password")
			assert.Falsef(t, utils.IsValidPassword("Mek!@#$%^"), "Expected week password")
		})

		t.Run("too short", func(t *testing.T) {
			t.Parallel()

			assert.Falsef(t, utils.IsValidPassword("Mek1!"), "Expected week password")
		})

		t.Run("too long", func(t *testing.T) {
			t.Parallel()

			too_long := "M9#dT7zW$k@1hXqE^6b!2pV3fU%4yG&sR*0jL8cQ+5nI?oA-CwDrZtHxuFgYKmPe!"
			assert.Falsef(t, utils.IsValidPassword(too_long), "Expected too long password")
		})
	})
}