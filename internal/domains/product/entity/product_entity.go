package entity

import (
	"fmt"
	"time"

	"github.com/royhairul/live-studio-api/internal/pkg/snowflakeid"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
	"gorm.io/gorm"
)

type Product struct {
	ID        string `gorm:"primaryKey"`
	UniqueID  string
	Name      string
	ShopID    string
	ShopName  string
	Platform  string
	Link      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	tenantdb.TenantBase
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if snowflakeid.Node == nil {
		return fmt.Errorf("snowflake node is not initialized")
	}
	p.ID = snowflakeid.Node.Generate().String()
	return nil
}
