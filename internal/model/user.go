package model

type User struct {
	Model
	Email    string `db:"email"`
	Password string `db:"password"`
}

type userUpdateField struct{}

func (userUpdateField) WithEmail(email string) UpdateField {
	return WithField("email", email)
}

func (userUpdateField) WithPassword(password string) UpdateField {
	return WithField("password", password)
}

var UserUpdate userUpdateField
