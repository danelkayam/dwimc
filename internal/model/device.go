package model

type Device struct {
	Model
	UserID      ID
	Serial      string
	Name        string
	Description string
	Token       string
}
