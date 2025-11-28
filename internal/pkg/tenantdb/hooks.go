package tenantdb

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// whitelist endpoint - these paths skip tenant filtering
var tenantWhitelist = []string{
	"/api/auth/login",
	"/api/auth/register",
	"/api/auth/forgot-password",
	"/api/auth/reset-password",
	"/api/auth/verify-otp",
}

func isWhitelisted(db *gorm.DB) bool {
	path := ExtractPath(db.Statement.Context)
	if path == "" {
		// no path available
		return false
	}

	for _, w := range tenantWhitelist {
		if strings.HasPrefix(path, w) {
			return true
		}
	}
	return false
}

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
	// Check whitelist first
	if isWhitelisted(db) {
		return
	}

	tenantID := ExtractTenant(db.Statement.Context)

	// Skip filtering if no tenant_id in context
	if tenantID == "" {
		return
	}

	stmt := db.Statement
	if stmt.Schema == nil {
		return
	}

	tableName := stmt.Schema.Table

	// ----------------------------------------------------
	// 1️⃣ Rule khusus: FILTER for TABLE "roles"
	// ----------------------------------------------------
	if tableName == "roles" {
		// Role global:       tenant_id IS NULL
		// Role per tenant:   tenant_id = current tenant
		db.Where("tenant_id = '' OR tenant_id = ?", tenantID)
		return
	}

	// ----------------------------------------------------
	// 2️⃣ Rule general: FILTER "tenant_id" for general if tenant_i exists
	// ----------------------------------------------------
	if _, ok := stmt.Schema.FieldsByName["TenantID"]; ok {
		db.Where("tenant_id = ?", tenantID)
	}
}
