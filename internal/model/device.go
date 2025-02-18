package model

type Device struct {
	Model
	UserID      ID
	Serial      string
	Name        string
	Description string
	Token       string
	// TODO - add last known location or all locations?
}
