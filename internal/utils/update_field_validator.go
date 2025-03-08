package utils

import "dwimc/internal/model"

type FieldValidator func(value any) error

type UpdateFieldsValidator struct {
	updateFields  []model.UpdateField
	validatorsMap map[string]FieldValidator
}

func NewUpdateFieldsValidator(updateFields []model.UpdateField) *UpdateFieldsValidator {
	return &UpdateFieldsValidator{
		updateFields:  updateFields,
		validatorsMap: map[string]FieldValidator{},
	}
}

func (ufv *UpdateFieldsValidator) WithValidator(field string, fieldValidator FieldValidator) *UpdateFieldsValidator {
	ufv.validatorsMap[field] = fieldValidator
	return ufv
}

func (ufv *UpdateFieldsValidator) Validate() error {
	fieldsValuesMap := map[string]any{}

	for _, updateField := range ufv.updateFields {
		updateField(&fieldsValuesMap)
	}

	for field, value := range fieldsValuesMap {
		if validator, ok := ufv.validatorsMap[field]; ok {
			if err := validator(value); err != nil {
				return err
			}
		}
	}

	return nil
}
