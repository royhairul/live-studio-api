package entity

import (
	"time"

	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
	"gorm.io/gorm"
)

type Shift struct {
	gorm.Model
	Name      string    `gorm:"type:varchar(100);not null;unique"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`

	tenantdb.TenantBase
}
