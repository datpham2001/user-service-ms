package passwordutil

import (
	"crypto/rand"
	"math/big"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

const (
	ResetPasswordTokenLength = 5
)

func HashPassword(pass string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateResetPasswordToken() (int, error) {
	tokenChars := make([]byte, ResetPasswordTokenLength)
	for i := range ResetPasswordTokenLength {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return 0, err
		}

		tokenChars[i] = byte(num.Int64() + '0')
	}

	token, err := strconv.Atoi(string(tokenChars))
	if err != nil {
		return 0, err
	}

	return token, nil
}
