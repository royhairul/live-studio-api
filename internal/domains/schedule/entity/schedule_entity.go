package entity

import (
	"time"

	"github.com/google/uuid"
	hostentity "github.com/royhairul/live-studio-api/internal/domains/host/entity"
	shiftentity "github.com/royhairul/live-studio-api/internal/domains/shift/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model

	HostID  uuid.UUID         `json:"host_id" gorm:"type:uuid;not null"`
	Host    hostentity.Host   `json:"host" gorm:"foreignKey:HostID"`
	ShiftID uint              `json:"shift_id" gorm:"not null"`
	Shift   shiftentity.Shift `json:"shift" gorm:"foreignKey:ShiftID"`

	Date      time.Time `json:"date" gorm:"not null"`
	StartTime time.Time `json:"start_time" gorm:"not null"`
	EndTime   time.Time `json:"end_time" gorm:"not null"`

	tenantdb.TenantBase
}
