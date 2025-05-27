package models

import (
	"time"

	"gorm.io/gorm"
)

type ResetPassword struct {
	gorm.Model
	Email string `gorm:"type:varchar(100)"`
	Otp string `gorm:"type:varchar(9);not null"`
	ExpiredAt time.Time
}