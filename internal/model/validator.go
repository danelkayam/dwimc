package model

type Validate func(value any) error

type Validator func(validates *map[string]Validate)

func WithValidator(field string, validate Validate) Validator {
	return func(validates *map[string]Validate) {
		(*validates)[field] = validate
	}
}
