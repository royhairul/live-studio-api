package entity

import (
	studioentity "github.com/royhairul/live-studio-api/internal/domains/studio/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name     string              `gorm:"type:varchar(100);not null"`
	UniqueID string              `gorm:"type:varchar(20);not null"`
	Username string              `gorm:"type:varchar(100);not null"`
	Password string              `gorm:"type:varchar(100)"`
	Email    string              `gorm:"type:varchar(100);not null"`
	Platform string              `gorm:"type:varchar(100);not null"`
	Cookie   string              `gorm:"type:text;not null"`
	Device   string              `gorm:"type:text"`
	StudioID uint16              `gorm:"not null"`
	Studio   studioentity.Studio `gorm:"foreignKey:StudioID;references:ID"`

	tenantdb.TenantBase
}
