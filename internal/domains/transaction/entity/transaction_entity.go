package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/pkg/snowflakeid"

	accountentity "github.com/royhairul/live-studio-api/internal/domains/account/entity"
	orderentity "github.com/royhairul/live-studio-api/internal/domains/order/entity"
)

type Transaction struct {
	ID                              int64  `gorm:"primaryKey"`
	UniqueID                        string `gorm:"unique"`
	Status                          string
	EstimatedTotalCommission        int64
	EstimatedTotalCommissionWithMCN int64

	PurchaseTime *time.Time `gorm:"null"`
	CompleteTime *time.Time `gorm:"null"`

	Orders []orderentity.Order `gorm:"foreignKey:TransactionID;references:ID"`

	AccountID uint
	Account   accountentity.Account `gorm:"foreignKey:AccountID;references:ID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (p *Transaction) BeforeCreate(tx *gorm.DB) error {
	if snowflakeid.Node == nil {
		return fmt.Errorf("snowflake node is not initialized")
	}
	p.ID = snowflakeid.Node.Generate().Int64()
	return nil
}
