package models

import (
	"gorm.io/gorm"
)

type Studio struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null"`
	Address string `gorm:"type:varchar(255)"`
}
