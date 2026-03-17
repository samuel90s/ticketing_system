package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var SECRET_KEY = []byte("secret123")

func GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(SECRET_KEY)
}
