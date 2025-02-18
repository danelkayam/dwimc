package model

import "time"

type ID uint

type Model struct {
	ID        			`db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
