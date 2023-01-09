package helper

import (
	"fmt"
	"os"
	"pendekin/dto"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func GenerateJWT(id uint, expiry time.Time) (string, error) {
	claim := dto.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiry.Unix(),
			Issuer:    os.Getenv("ISSUER"),
		},
		UserId: id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	jwt_secret_key := os.Getenv("JWT_SECRET_KEY")
	if jwt_secret_key == "" {
		jwt_secret_key = "default-short-jwt-secret-key-123"
	}

	signedToken, err := token.SignedString([]byte(jwt_secret_key))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(token string) (*jwt.Token, error) {
	jwt_secret_key := os.Getenv("JWT_SECRET_KEY")
	if jwt_secret_key == "" {
		jwt_secret_key = "default-short-jwt-secret-key-123"
	}

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Token is invalid!")
		}
		return []byte(jwt_secret_key), nil
	})
}
