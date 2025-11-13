package tenantdb

import "context"

type tenantKey string

const TenantKey tenantKey = "tenant_id"

var globalTenant string

// SetTenant global
func SetTenant(id string) {
	globalTenant = id
}

// GetTenant
func GetTenant() string {
	return globalTenant
}

// AttachTenant
func AttachTenant(ctx context.Context, tenantID string) context.Context {
	SetTenant(tenantID)
	return context.WithValue(ctx, TenantKey, tenantID)
}

// ExtractTenant
func ExtractTenant(ctx context.Context) string {
	if v := ctx.Value(TenantKey); v != nil {
		if tenant, ok := v.(string); ok {
			return tenant
		}
	}
	return globalTenant
}
