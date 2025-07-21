package entity

import (
	"gorm.io/gorm"

	permissionentity "github.com/royhairul/live-studio-api/internal/domains/permission/entity"
)

type Role struct {
	gorm.Model
	Name        string                        `json:"name"`
	Permissions []permissionentity.Permission `json:"permissions" gorm:"many2many:role_permissions"`
}
