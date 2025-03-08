package model

type UpdateField func(updateFields *map[string]any)

func WithField(field string, value any) UpdateField {
	return func(updateFields *map[string]any) {
		(*updateFields)[field] = value
	}
}