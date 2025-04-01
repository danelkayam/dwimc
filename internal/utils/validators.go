package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
	RegisterValidations(validate)
}

func GetDefaultValidate() *validator.Validate {
	return validate
}

func RegisterValidations(v *validator.Validate) {
	err := v.RegisterValidation("nonempty", nonempty)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to register validator: nonempty")
	}
}

// Custom validates definitions

func nonempty(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(field.String())) > 0

	case reflect.Slice, reflect.Array, reflect.Map:
		return field.Len() > 0

	default:
		return false
	}
}
