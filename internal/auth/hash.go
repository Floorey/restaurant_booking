package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(pw string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pw), 12)
}

func Verify(hash []byte, pw string) bool {
	return bcrypt.CompareHashAndPassword(hash, []byte(pw)) == nil
}
