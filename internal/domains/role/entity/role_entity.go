package entity

import (
	permissionentity "github.com/royhairul/live-studio-api/internal/domains/permission/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string
	Permissions []permissionentity.Permission `gorm:"many2many:role_permissions"`

	tenantdb.TenantBase
}
