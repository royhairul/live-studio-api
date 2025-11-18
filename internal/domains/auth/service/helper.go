package service

import (
	"github.com/royhairul/live-studio-api/database"
	permissionentity "github.com/royhairul/live-studio-api/internal/domains/permission/entity"
	roleentity "github.com/royhairul/live-studio-api/internal/domains/role/entity"
	"github.com/royhairul/live-studio-api/internal/pkg/tenantdb"
)

func CreateTenantRoles(tenantID string) ([]roleentity.Role, error) {
	roles := []struct {
		Name        string
		Permissions []string
	}{
		{
			Name: "admin",
			Permissions: []string{
				"host_view", "host_create", "host_edit", "host_delete",
				"schedule_host_view", "schedule_host_create", "schedule_host_edit", "schedule_host_delete",
				"shift_view", "shift_create", "shift_edit", "shift_delete",
				"account_view", "account_create", "account_edit", "account_delete",
				"live_view", "live_report",
				"finance_view", "finance_research_view", "finance_research_create",
			},
		},
		{
			Name: "host",
			Permissions: []string{
				"schedule_host_view",
				"shift_view",
				"live_view",
			},
		},
	}

	var createdRoles []roleentity.Role

	for _, r := range roles {
		role := roleentity.Role{
			Name: r.Name,
			TenantBase: tenantdb.TenantBase{
				TenantID: tenantID,
			},
		}

		database.DB.Create(&role)

		// Hubungkan permission
		var perms []permissionentity.Permission
		database.DB.Where("name IN ?", r.Permissions).Find(&perms)

		database.DB.Model(&role).Association("Permissions").Replace(perms)

		createdRoles = append(createdRoles, role)
	}

	return createdRoles, nil
}
