package entity

import (
	"time"

	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	"gorm.io/gorm"
)

type Accountads struct {
	gorm.Model

	Spend uint       `gorm:"not null"`
	Date  *time.Time `gorm:"index;not null"`

	// Relasi ke akun
	AccountID uint
	Account   accountentity.Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
