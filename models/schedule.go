package models

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	HostID  uint          `json:"host_id" gorm:"not null"`
	Host    Host          `json:"host" gorm:"foreignKey:HostID"`
	ShiftID uint          `json:"shift_id" gorm:"not null"`
	Shift   ScheduleShift `json:"shift" gorm:"foreignKey:ShiftID"`

	Date      time.Time `json:"date" gorm:"not null"`
	StartTime time.Time `json:"start_time" gorm:"not null"`
	EndTime   time.Time `json:"end_time" gorm:"not null"`
}
