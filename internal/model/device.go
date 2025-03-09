package model

import "database/sql"

type Device struct {
	Model
	UserID ID             `db:"user_id"`
	Serial string         `db:"serial"`
	Name   string         `db:"name"`
	Token  sql.NullString `db:"token"`
}

func WithSerial(serial string) Field {
	return WithField("serial", serial)
}

func WithName(name string) Field {
	return WithField("name", name)
}

func WithToken(token string) Field {
	return WithField("token", token)
}
