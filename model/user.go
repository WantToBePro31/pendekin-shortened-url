package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Fullname string `json:"fullname" gorm:"not null"`
	Username string `json:"username" gorm:"unique;not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
}
