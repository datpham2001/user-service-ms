package tokensvc

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	AccessTokenDuration  = time.Hour * 24
	RefreshTokenDuration = time.Hour * 24 * 30
)

type JwtToken struct {
	secretKey string
}

func NewJwtToken(secretKey string) *JwtToken {
	return &JwtToken{
		secretKey: secretKey,
	}
}

func (t *JwtToken) GenerateTokenPair(userId string) (string, string, error) {
	accessToken, err := t.generateToken(userId, AccessTokenDuration)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := t.generateToken(userId, RefreshTokenDuration)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (t *JwtToken) generateToken(userId string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(duration).Unix(),
	})

	return token.SignedString([]byte(t.secretKey))
}
