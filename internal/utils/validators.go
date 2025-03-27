package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	err := validate.RegisterValidation("nonempty", nonempty)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to register validator: nonempty")
	}
}

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
