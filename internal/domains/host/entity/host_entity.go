package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	studioentity "github.com/royhairul/live-studio-api/internal/domains/studio/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
)

type Host struct {
	ID     *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name   string     `gorm:"type:varchar(100);not null"`
	Phone  string     `gorm:"type:varchar(15);not null"`
	UserID uint

	StudioID uint                `gorm:"not null"`
	Studio   studioentity.Studio `gorm:"foreignKey:StudioID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	tenantdb.TenantBase
}
