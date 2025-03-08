package model

import "database/sql"

type User struct {
	Model
	Email    string         `db:"email"`
	Password string         `db:"password"`
	Token    sql.NullString `db:"token"`
}

type userUpdateField struct{}

func (userUpdateField) WithPassword(password string) UpdateField {
	return WithField("password", password)
}

func (userUpdateField) WithToken(token string) UpdateField {
	return WithField("token", token)
}

var UserUpdate userUpdateField
