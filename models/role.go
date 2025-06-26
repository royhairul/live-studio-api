package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions"`
}
