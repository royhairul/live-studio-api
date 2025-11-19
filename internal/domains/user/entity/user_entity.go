package entity

import (
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	roleentity "github.com/royhairul/live-studio-api/internal/domains/role/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
)

type User struct {
	ID       *uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name     string          `gorm:"type:varchar(100);not null"`
	Email    string          `gorm:"type:varchar(100)"`
	Password string          `gorm:"type:varchar(100);not null"`
	RoleID   uint            `gorm:"not null"`
	Role     roleentity.Role `gorm:"foreignKey:RoleID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	tenantdb.TenantBase
}
