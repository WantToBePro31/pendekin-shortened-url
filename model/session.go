package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	JWT            string    `json:"jwt" gorm:"unique;not null"`
	ExpirationTime time.Time `json:"expiration_time" gorm:"not null"`
	UserId         uint      `json:"user_id" gorm:"not null"`
}
