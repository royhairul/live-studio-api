package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	hostentity "github.com/royhairul/live-studio-api/internal/domains/host/entity"
	scheduleentity "github.com/royhairul/live-studio-api/internal/domains/schedule/entity"
	shiftentity "github.com/royhairul/live-studio-api/internal/domains/shift/entity"
)

type Attendance struct {
	gorm.Model
	ScheduleID *uint                   `json:"schedule_id,omitempty"`
	Schedule   scheduleentity.Schedule `json:"schedule" gorm:"foreignKey:ScheduleID"`

	HostID *uuid.UUID      `json:"host_id"`
	Host   hostentity.Host `json:"host" gorm:"foreignKey:HostID"`

	ShiftID uint              `json:"shift_id"`
	Shift   shiftentity.Shift `json:"shift" gorm:"foreignKey:ShiftID"`

	Date         *time.Time `json:"date"`
	CheckedInAt  *time.Time `json:"checked_in,omitempty"`
	CheckedOutAt *time.Time `json:"checked_out,omitempty"`

	Status string `json:"status"`
	Note   string `json:"note"`
}
