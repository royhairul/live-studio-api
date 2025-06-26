package models

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	gorm.Model
	ScheduleID uint     `json:"schedule_id"`
	Schedule   Schedule `json:"schedule"`

	CheckedInAt  *time.Time `json:"checked_in"`
	CheckedOutAt *time.Time `json:"checked_out"`

	Status string
	Note   string
}
