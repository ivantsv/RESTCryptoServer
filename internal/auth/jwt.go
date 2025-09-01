package auth

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrInvalidToken = errors.New("invalid token")  

func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY_CRYPTO_SERVER")))
	if err != nil {
		log.Println("error during token creation")
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidToken
        }
		return []byte(os.Getenv("SECRET_KEY_CRYPTO_SERVER")), nil
	})
	if err != nil {
		log.Println("error during token validation")
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token, nil
}