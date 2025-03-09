package model

type Field func(fields *map[string]any)

func WithField(field string, value any) Field {
	return func(fields *map[string]any) {
		(*fields)[field] = value
	}
}