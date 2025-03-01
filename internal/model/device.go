package model

import "database/sql"

type Device struct {
	Model
	UserID ID             `db:"user_id"`
	Serial string         `db:"serial"`
	Name   string         `db:"name"`
	Token  sql.NullString `db:"token"`
}
