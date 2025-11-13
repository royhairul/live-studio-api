package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	hostentity "github.com/royhairul/live-studio-api/internal/domains/host/entity"
	scheduleentity "github.com/royhairul/live-studio-api/internal/domains/schedule/entity"
	shiftentity "github.com/royhairul/live-studio-api/internal/domains/shift/entity"
	studioentity "github.com/royhairul/live-studio-api/internal/domains/studio/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
)

type Attendance struct {
	gorm.Model
	Date         *time.Time
	CheckedInAt  *time.Time
	CheckedOutAt *time.Time

	Status string
	Note   string

	ScheduleID *uint
	Schedule   scheduleentity.Schedule `gorm:"foreignKey:ScheduleID"`

	HostID *uuid.UUID      `gorm:"type:uuid"`
	Host   hostentity.Host `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	ShiftID uint
	Shift   shiftentity.Shift `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	StudioID uint
	Studio   studioentity.Studio `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	tenantdb.TenantBase
}

func (a *Attendance) Duration() int64 {
	if a.CheckedInAt != nil && a.CheckedOutAt != nil {
		return int64(a.CheckedOutAt.Sub(*a.CheckedInAt).Seconds())
	}
	return 0
}
