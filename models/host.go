package models

import (
	"gorm.io/gorm"
)

type Host struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100);not null"`
	Phone    string `gorm:"type:varchar(15);not null"`
	StudioID uint16 `gorm:"not null"`
	Studio   Studio `gorm:"foreignKey:StudioID;references:ID"`
	UserID   uint16 `gorm:"not null"`
}
