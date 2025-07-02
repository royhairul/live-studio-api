package entity

import (
	"time"

	"gorm.io/gorm"
)

type Shift struct {
	gorm.Model
	Name      string    `gorm:"type:varchar(100);not null;unique"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
}
