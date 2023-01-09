package dto

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	jwt.StandardClaims
	UserId uint `json:"user_id"`
}

type Session struct {
	JWT            string    `json:"jwt" binding:"required"`
	ExpirationTime time.Time `json:"expiration_time" binding:"required"`
	UserId         uint      `json:"user_id" binding:"required"`
}
