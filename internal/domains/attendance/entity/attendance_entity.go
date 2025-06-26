package entity

import (
	"time"

	"github.com/royhairul/live-studio-api/models"
	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	ScheduleID *uint            `json:"schedule_id,omitempty"`
	Schedule   *models.Schedule `json:"schedule" gorm:"foreignKey:ScheduleID"`

	ShiftID *uint                 `json:"shift_id"`
	Shift   *models.ScheduleShift `json:"shift" gorm:"foreignKey:ShiftID"`

	Date         *time.Time `json:"date"`
	CheckedInAt  *time.Time `json:"checked_in,omitempty"`
	CheckedOutAt *time.Time `json:"checked_out,omitempty"`

	HostID *uint       `json:"host_id"`
	Host   models.Host `json:"host" gorm:"foreignKey:HostID"`

	Status string `json:"status"`
	Note   string `json:"note"`
}
