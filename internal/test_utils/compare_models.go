package testutils

import (
	"reflect"
)

type OnFailedCheck func(field string, shouldBeEqual bool, got any, expected any)

type CheckField[T any] func(checkFields *map[string]bool)

func WithFieldNotEqual[T any](fieldName string) CheckField[T] {
	return func(checkFields *map[string]bool) {
		(*checkFields)[fieldName] = false
	}
}

func AssertEqualItems[T any](expected *T, got *T, onFailedCheck OnFailedCheck, opts ...CheckField[T]) {
	checkFields := map[string]bool{}

	expectedType := reflect.TypeOf(expected).Elem()

	setDefaultCheckFields(expectedType, &checkFields)

	for _, opt := range opts {
		opt(&checkFields)
	}

	expectedValue := reflect.ValueOf(expected).Elem()
	gotValue := reflect.ValueOf(got).Elem()

	compareFields(expectedValue, gotValue, &checkFields, &onFailedCheck)
}

func setDefaultCheckFields(structType reflect.Type, checkFields *map[string]bool) {
	for i := range structType.NumField() {
		field := structType.Field(i)

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			setDefaultCheckFields(field.Type, checkFields)
		} else {
			(*checkFields)[field.Name] = true
		}
	}
}

func compareFields(expected reflect.Value, got reflect.Value, checkFields *map[string]bool, onFailedCheck *OnFailedCheck) {
	for i := range expected.NumField() {
		field := expected.Type().Field(i)
		expectedField := expected.Field(i)
		gotField := got.Field(i)

		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			compareFields(expectedField, gotField, checkFields, onFailedCheck)
			continue
		}

		equality := (*checkFields)[field.Name]
		if equality {
			if !reflect.DeepEqual(gotField.Interface(), expectedField.Interface()) {
				(*onFailedCheck)(field.Name, equality, gotField.Interface(), expectedField.Interface())
			}
		} else {
			if reflect.DeepEqual(gotField.Interface(), expectedField.Interface()) {
				(*onFailedCheck)(field.Name, equality, gotField.Interface(), expectedField.Interface())
			}
		}
	}
}
