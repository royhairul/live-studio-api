package entity

import (
	models "github.com/royhairul/live-studio-api/models"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name     string        `gorm:"type:varchar(100);not null"`
	UniqueID string        `gorm:"type:varchar(20);not null;unique"`
	Username string        `gorm:"type:varchar(100);not null"`
	Password string        `gorm:"type:varchar(100)"`
	Email    string        `gorm:"type:varchar(100);not null"`
	Platform string        `gorm:"type:varchar(100);not null"`
	Cookie   string        `gorm:"type:text;not null"`
	Device   string        `gorm:"type:text"`
	StudioID uint16        `gorm:"not null"`
	Studio   models.Studio `gorm:"foreignKey:StudioID;references:ID"`
}
