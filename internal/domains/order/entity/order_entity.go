package entity

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/helpers/snowflakeid"
	productentity "github.com/royhairul/live-studio-api/internal/domains/product/entity"
)

type Order struct {
	ID            string `gorm:"primaryKey"`
	SerialNumber  string
	Status        string
	CompleteTime  int64
	TransactionID int64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	Products []productentity.Product `gorm:"many2many:order_items;"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if snowflakeid.Node == nil {
		return fmt.Errorf("snowflake node is not initialized")
	}
	o.ID = snowflakeid.Node.Generate().String()
	return nil
}
