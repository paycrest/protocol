package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash is a function to compare provided password with the hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
