package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password string, hashedPassword string) bool {
	// Compare plain password with stored hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}