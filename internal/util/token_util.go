package util

import (
	"crypto/rand"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	ResetPasswordTokenLength = 5
)

func GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
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
