package models

import (
	"gorm.io/gorm"
)

type Studio struct {
	gorm.Model
	Nomor string `gorm:"type:varchar(100);not null"`
}