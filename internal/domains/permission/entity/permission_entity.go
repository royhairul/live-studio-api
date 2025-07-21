package entity

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique"`
	Group       string `json:"group"`
	Description string `json:"description"`
}
