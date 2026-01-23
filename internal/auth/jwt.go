package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateToken(userId uuid.UUID, secret string) (string, error) {
	expiresAt := time.Now().Add(48 * time.Hour)

	claims := jwt.MapClaims{
		"sub": userId.String(),
		"exp": expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}
func ValidateToken(tokenStr, secret string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("error occurred")
	}
	subStr, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, errors.New("sub is not string")
	}

	id, err := uuid.Parse(subStr)

	return id, err
}
