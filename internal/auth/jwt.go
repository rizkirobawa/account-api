package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret")

var ErrInvalidToken = errors.New("token tidak valid")

// Claims adalah data yang "dibawa" oleh token
type Claims struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken membuat token baru untuk satu akun, berlaku 24 jam
func GenerateToken(accountID, email string) (string, error) {
	claims := Claims{
		AccountID: accountID,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
