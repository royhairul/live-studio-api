package entity

import (
	"time"

	"gorm.io/gorm"
)

type ResetPassword struct {
	gorm.Model
	Email     string `gorm:"unique,type:varchar(255)"`
	Otp       string `gorm:"unique,type:varchar(255)"`
	ExpiredAt time.Time
}
