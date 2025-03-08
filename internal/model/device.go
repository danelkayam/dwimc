package model

import "database/sql"

type Device struct {
	Model
	UserID ID             `db:"user_id"`
	Serial string         `db:"serial"`
	Name   string         `db:"name"`
	Token  sql.NullString `db:"token"`
}

type deviceUpdateField struct{}

func (deviceUpdateField) WithSerial(serial string) UpdateField {
	return WithField("serial", serial)
}

func (deviceUpdateField) WithName(name string) UpdateField {
	return WithField("name", name)
}

func (deviceUpdateField) WithToken(token string) UpdateField {
	return WithField("token", token)
}

var DeviceUpdate deviceUpdateField
