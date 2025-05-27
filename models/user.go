package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);not null"`
	Email	 string `gorm:"type:varchar(100)"`
	Password string `gorm:"type:varchar(100);not null"`
	Role	 string	`gorm:"type:varchar(10);not null;default:'host'"`
}