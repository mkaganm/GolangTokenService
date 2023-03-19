package auth

import "golang.org/x/crypto/bcrypt"

// PasswordToHash This function converts password to hash.
func PasswordToHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword This function returns a boolean value
// by comparing the password and hash.
func CheckPassword(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
