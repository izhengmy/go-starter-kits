package crypt

import "golang.org/x/crypto/bcrypt"

func PasswordHash(value string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func PasswordVerify(value string, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(value)) == nil
}
