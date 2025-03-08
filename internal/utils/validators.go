package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	err := validate.RegisterValidation("strong_password", PasswordValidator)
	if err != nil {
		// WTF moment - should not get here
		panic(err)
	}
}

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	var (
		hasMinLen  = len(password) >= 8
		hasMaxLen  = len(password) <= 64
		hasNumber  = regexp.MustCompile(`[0-9]`).MatchString(password)
		hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
		hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
		hasSpecial = regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
	)
	return hasMinLen && hasMaxLen && hasNumber && hasUpper && hasLower && hasSpecial
}

func IsValidEmail(email string) bool {
	return validate.Var(email, "required,email,min=5,max=254") == nil
}

func IsValidPassword(password string) bool {
	return validate.Var(password, "required,min=8,max=64,strong_password") == nil
}
