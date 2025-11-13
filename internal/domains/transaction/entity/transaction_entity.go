package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/pkg/snowflakeid"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"

	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	orderentity "github.com/royhairul/live-studio-api/internal/domains/order/entity"
)

type Transaction struct {
	ID                     int64  `gorm:"primaryKey"`
	UniqueID               string `gorm:"unique"`
	Status                 string
	TotalPurchase          int64
	TotalCommission        int64
	TotalCommissionWithMCN int64
	PurchaseTime           *time.Time
	CompleteTime           *time.Time
	Orders                 []orderentity.Order `gorm:"foreignKey:TransactionID;references:ID"`
	AccountID              uint
	Account                accountentity.Account `gorm:"foreignKey:AccountID;references:ID;preload:Studio"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt `gorm:"index"`

	tenantdb.TenantBase
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if snowflakeid.Node == nil {
		return fmt.Errorf("snowflake node is not initialized")
	}
	t.ID = snowflakeid.Node.Generate().Int64()
	return nil
}
