package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the plainTextPassword
func HashPassword(plainTextPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("falied to hash password, error:%s", err)
	}
	return string(hashedPassword), nil
}
