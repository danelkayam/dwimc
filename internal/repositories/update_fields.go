package repositories

type UpdateField func(updateFields *map[string]interface{})

func WithField(field string, value interface{}) UpdateField {
	return func(updateFields *map[string]interface{}) {
		(*updateFields)[field] = value
	}
}