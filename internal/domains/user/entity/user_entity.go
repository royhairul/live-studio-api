package entity

import (
	"gorm.io/gorm"

	roleentity "github.com/royhairul/live-studio-api/internal/domains/role/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
)

type User struct {
	gorm.Model
	Name     string          `gorm:"type:varchar(100);not null"`
	Email    string          `gorm:"type:varchar(100)"`
	Password string          `gorm:"type:varchar(100);not null"`
	RoleID   uint            `gorm:"not null"`
	Role     roleentity.Role `gorm:"foreignKey:RoleID;references:ID"`

	tenantdb.TenantBase
}
