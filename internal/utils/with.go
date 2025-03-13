package utils

import "dwimc/internal/model"

type Validate func(value any) error

type Validator func(validates *map[string]Validate)

func WithValidator(field string, validate Validate) Validator {
	return func(validates *map[string]Validate) {
		(*validates)[field] = validate
	}
}

type With struct {
	fieldsValuesMap    map[string]any
	fieldsValidatesMap map[string]Validate
}

func NewWithValidator() *With {
	return &With{
		fieldsValuesMap:    map[string]any{},
		fieldsValidatesMap: map[string]Validate{},
	}
}

func (with *With) WithField(field model.Field) *With {
	field(&with.fieldsValuesMap)
	return with
}

func (with *With) WithFields(fields []model.Field) *With {
	for _, field := range fields {
		field(&with.fieldsValuesMap)
	}
	return with
}

func (with *With) WithValidator(validator Validator) *With {
	validator(&with.fieldsValidatesMap)
	return with
}

func (with *With) WithValidators(validators []Validator) *With {
	for _, validator := range validators {
		validator(&with.fieldsValidatesMap)
	}
	return with
}

func (with *With) Validate() error {
	for field, value := range with.fieldsValuesMap {
		if validate, ok := with.fieldsValidatesMap[field]; ok {
			if err := validate(value); err != nil {
				return err
			}
		}
	}

	return nil
}
