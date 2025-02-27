package model

type Device struct {
	Model
	UserID      ID
	Serial      string
	Name        string
	Token       string
}
