package entity

import (
	"time"

	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
	"gorm.io/gorm"
)

type Accountads struct {
	gorm.Model

	Spend uint       `gorm:"not null"`
	Date  *time.Time `gorm:"uniqueIndex:uidx_account_date,not null"`

	// Relasi ke akun
	AccountID uint                  `gorm:"uniqueIndex:uidx_account_date"`
	Account   accountentity.Account `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	tenantdb.TenantBase
}
