package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func AsError(err error, reason string) error {
	return fmt.Errorf("%s: %w", reason, err)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
