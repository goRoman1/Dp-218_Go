package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("cant hash user password %w", err)
	}
	hash := string(hashedPassword)
	return hash, nil
}

func CheckPassword(password1, password2 string) error {
	return bcrypt.CompareHashAndPassword([]byte(password1), []byte(password2))
}

// func Sanitize() {
// 	u.Password = ""
// 	u.Role = ""
// }
