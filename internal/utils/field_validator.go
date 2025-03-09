package utils

import "dwimc/internal/model"

type FieldValidator func(value any) error

type FieldsValidator struct {
	fields        []model.Field
	validatorsMap map[string]FieldValidator
}

func NewFieldsValidator() *FieldsValidator {
	return &FieldsValidator{
		fields:        []model.Field{},
		validatorsMap: map[string]FieldValidator{},
	}
}

func (fv *FieldsValidator) WithFields(fields []model.Field) *FieldsValidator {
	fv.fields = append(fv.fields, fields...)
	return fv
}

func (fv *FieldsValidator) WithField(field model.Field) *FieldsValidator {
	fv.fields = append(fv.fields, field)
	return fv
}

func (fv *FieldsValidator) WithValidator(field string, fieldValidator FieldValidator) *FieldsValidator {
	fv.validatorsMap[field] = fieldValidator
	return fv
}


func (fv *FieldsValidator) Validate() error {
	fieldsValuesMap := map[string]any{}

	for _, field := range fv.fields {
		field(&fieldsValuesMap)
	}

	for field, value := range fieldsValuesMap {
		if validator, ok := fv.validatorsMap[field]; ok {
			if err := validator(value); err != nil {
				return err
			}
		}
	}

	return nil
}
