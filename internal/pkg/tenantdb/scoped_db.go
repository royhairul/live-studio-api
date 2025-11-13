package tenantdb

import (
	"context"

	"gorm.io/gorm"
)

type TenantScope struct {
	DB *gorm.DB
}

func WithTenant(db *gorm.DB, ctx context.Context) *gorm.DB {
	tenantID, ok := ctx.Value("tenant_id").(string)
	if !ok || tenantID == "" {
		return db
	}
	return db.Scopes(func(tx *gorm.DB) *gorm.DB {
		return tx.Where("tenant_id = ?", tenantID)
	})
}
