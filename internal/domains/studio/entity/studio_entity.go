package entity

import (
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
	"gorm.io/gorm"
)

type Studio struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);not null"`
	Address string `gorm:"type:varchar(255)"`

	tenantdb.TenantBase
}
