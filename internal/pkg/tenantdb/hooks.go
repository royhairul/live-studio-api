package tenantdb

import (
	"reflect"

	"gorm.io/gorm"
)

// RegisterTenantCallback sets tenant_id automatically for tenant models
func RegisterTenantCallback(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").Register("tenantdb:set_tenant_id", setTenantID)
	db.Callback().Query().Before("gorm:query").Register("tenantdb:filter_tenant_id", filterByTenantID)
}

func setTenantID(db *gorm.DB) {
	tenantID := ExtractTenant(db.Statement.Context)
	if tenantID == "" {
		return
	}

	rv := reflect.Indirect(reflect.ValueOf(db.Statement.Dest))
	if rv.Kind() == reflect.Struct {
		field := rv.FieldByName("TenantID")
		if field.IsValid() && field.CanSet() && field.Kind() == reflect.String {
			field.SetString(tenantID)
		}
	}
}

func filterByTenantID(db *gorm.DB) {
	tenantID := ExtractTenant(db.Statement.Context)
	if tenantID == "" {
		return
	}

	stmt := db.Statement
	if stmt.Schema != nil {
		if _, ok := stmt.Schema.FieldsByName["TenantID"]; ok {
			db.Where("tenant_id = ?", tenantID)
		}
	}
}
