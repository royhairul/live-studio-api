package params

import (
	permissionparams "github.com/royhairul/live-studio-api/internal/domains/permission/params"
	"github.com/royhairul/live-studio-api/internal/domains/role/entity"
)

type RoleResponse struct {
	// TODO: add response fields
	ID          uint                                  `json:"id"`
	Name        string                                `json:"name"`
	Permissions []permissionparams.PermissionResponse `json:"permissions"`
}

func NewRoleResponse(role *entity.Role) *RoleResponse {
	var permissions []permissionparams.PermissionResponse
	for _, permission := range role.Permissions {
		permissions = append(permissions, *permissionparams.NewPermissionResponse(&permission))
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
	}
}
