package params

import (
	"github.com/royhairul/live-studio-api/internal/domains/role/entity"
)

type RoleResponse struct {
	// TODO: add response fields
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func NewRoleResponse(role *entity.Role) *RoleResponse {
	permissions := make([]string, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		permissions = append(permissions, p.Name)
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissions,
	}
}
