package utils

import "dwimc/internal/model"

type With struct {
	fieldsValuesMap    map[string]any
	fieldsValidatesMap map[string]model.Validate
}

func NewWithValidator() *With {
	return &With{}
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

func (with *With) WithValidator(validator model.Validator) *With {
	validator(&with.fieldsValidatesMap)
	return with
}

func (with *With) WithValidators(validators []model.Validator) *With {
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
