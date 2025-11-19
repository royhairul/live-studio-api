package tenantdb

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TenantBase struct {
	TenantID string `json:"tenant_id" gorm:"index;nullable"`
}

func (t *TenantBase) BeforeCreate(tx *gorm.DB) (err error) {
	if t.TenantID == "" {
		if ctxTenant, ok := tx.Statement.Context.Value("tenant_id").(string); ok && ctxTenant != "" {
			t.TenantID = ctxTenant
		}
		// else {
		// 	// fallback jika tidak ada context (opsional)
		// 	t.TenantID = uuid.New().String()
		// }
	}
	return
}

func (t *TenantBase) BeforeFind(tx *gorm.DB) (err error) {
	if ctxTenant, ok := tx.Statement.Context.Value("tenant_id").(string); ok && ctxTenant != "" {
		tx.Statement.AddClause(clause.Where{
			Exprs: []clause.Expression{
				clause.Eq{Column: "tenant_id", Value: ctxTenant},
			},
		})
	}
	return
}
