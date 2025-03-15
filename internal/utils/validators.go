package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("strong_password", strongPassword)
	if err != nil {
		// WTF moment - should not get here
		panic(err)
	}
}

func strongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var (
		hasMinLen  = len(password) >= 8
		hasMaxLen  = len(password) <= 64
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
		hasNot     = regexp.MustCompile(`[^\s\t\;]`).MatchString(password)
	)
	return hasMinLen && hasMaxLen && hasNumber && hasUpper && hasLower && hasSpecial && hasNot
}
