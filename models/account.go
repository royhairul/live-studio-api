package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null"`
	UniqueID string `gorm:"type:varchar(20);not null;unique"`
	Username string `gorm:"type:varchar(100);not null"`
	Password string `gorm:"type:varchar(100)"`
	Email    string `gorm:"type:varchar(100);not null"`
	Platform string `gorm:"type:varchar(100);not null"`
	Cookies  string `gorm:"type:text;not null"`
	StudioID uint16 `gorm:"not null"`
	Studio   Studio `gorm:"foreignKey:StudioID;references:ID"`
}
