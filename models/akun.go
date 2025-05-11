package models

import (
	"gorm.io/gorm"
)

type Akun struct {
	gorm.Model
	Name	 string `gorm:"type:varchar(100);not null"`
	Username string `gorm:"type:varchar(100);not null"`
	Password string `gorm:"type:varchar(100)"`
	Email	 string `gorm:"type:varchar(100);not null;unique"`
	Platform  string `gorm:"type:varchar(100);not null"`
	Cookies	 string `gorm:"type:text;not null"`
}