package tenantdb

import "context"

type tenantKey string

const TenantKey tenantKey = "tenant_id"

var (
	globalTenant string
	globalPath   string
)

// SetTenant global
func SetTenant(id string) {
	globalTenant = id
}

// SetRequestPath sets a global request path (fallback when context path not available)
func SetRequestPath(path string) {
	globalPath = path
}

// GetRequestPath returns the global request path
func GetRequestPath() string {
	return globalPath
}

// GetTenant
func GetTenant() string {
	return globalTenant
}

// AttachTenant
func AttachTenant(ctx context.Context, tenantID string) context.Context {
	// do NOT set global tenant here to avoid leaking tenant across requests
	return context.WithValue(ctx, TenantKey, tenantID)
}

// ExtractTenant
func ExtractTenant(ctx context.Context) string {
	if v := ctx.Value(TenantKey); v != nil {
		if tenant, ok := v.(string); ok {
			return tenant
		}
	}
	// Do not fall back to global tenant to prevent cross-request leakage
	return ""
}

// ExtractPath returns path from context or falls back to globalPath
func ExtractPath(ctx context.Context) string {
	if v := ctx.Value("path"); v != nil {
		if p, ok := v.(string); ok && p != "" {
			return p
		}
	}
	return globalPath
}
