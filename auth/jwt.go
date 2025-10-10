package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var key string = os.Getenv("JWT_SECRET")

func GenerateAccessToken(userId int64, username string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":     "my-todo-app",
		"sub":     username,
		"user_id": userId,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})
	s, err := t.SignedString([]byte(key))
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}
	return s, nil
}

func VerifyAccessToken(token string) (jwt.MapClaims, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(key), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	if t.Valid {
		return t.Claims.(jwt.MapClaims), nil
	}

	return nil, fmt.Errorf("invalid token: %w", err)
}
