package models

import "gorm.io/gorm"

type (
	User struct {
		gorm.Model
		Email    string `gorm:"unique;not null"`
		Username string `gorm:"unique;not null"`
		Password string `gorm:"not null"`
	}
	SignReq struct {
		Email    string
		Username string
		Password string
	}
)
