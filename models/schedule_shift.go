package models

import (
	"time"

	"gorm.io/gorm"
)

type ScheduleShift struct {
	gorm.Model
	Name      string    `json:"name" gorm:"type:varchar(100);not null;unique"`
	StartTime time.Time `json:"start_time" gorm:"not null"`
	EndTime   time.Time `json:"end_time" gorm:"not null"`
}
