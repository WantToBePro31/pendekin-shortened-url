package model

import "gorm.io/gorm"

type Url struct {
	gorm.Model
	RealUrl      string `json:"real_url" gorm:"not null"`
	ShortenedUrl string `json:"shortened_url" gorm:"unique;not null"`
	Randomized   bool   `json:"randomized" gorm:"not null"`
	UserId       uint   `json:"user_id" gorm:"not null"`
}
