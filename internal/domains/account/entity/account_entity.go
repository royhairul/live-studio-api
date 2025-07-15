package entity

import (
	models "github.com/royhairul/live-studio-api/models"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name     string        `gorm:"type:varchar(100);not null" json:"name"`
	UniqueID string        `gorm:"type:varchar(20);not null;unique" json:"unique_id"`
	Username string        `gorm:"type:varchar(100);not null" json:"username"`
	Password string        `gorm:"type:varchar(100)" json:"-"`
	Email    string        `gorm:"type:varchar(100);not null" json:"email"`
	Platform string        `gorm:"type:varchar(100);not null" json:"platform"`
	Cookie   string        `gorm:"type:text;not null" json:"cookie"`
	StudioID uint16        `gorm:"not null" json:"studio_id"`
	Studio   models.Studio `gorm:"foreignKey:StudioID;references:ID" json:"studio"`
}
