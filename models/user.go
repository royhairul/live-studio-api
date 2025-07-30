package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null"`
	Email    string `gorm:"type:varchar(100)"`
	Password string `gorm:"type:varchar(100);not null"`
	RoleID   uint   `gorm:"not null"`
	Role     Role   `gorm:"foreignKey:RoleID;references:ID`
}
