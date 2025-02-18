package model

import "database/sql"

type User struct {
	Model
	Email    string         `db:"email"`
	Password string         `db:"password"`
	Token    sql.NullString `db:"token"`
}
