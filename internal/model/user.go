package model

type User struct {
	Model
	Email    string `db:"email"`
	Password string `db:"password"`
}

func WithEmail(email string) Field {
	return WithField("email", email)
}

func WithPassword(password string) Field {
	return WithField("password", password)
}
