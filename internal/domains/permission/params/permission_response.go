package params

import "github.com/royhairul/live-studio-api/internal/domains/permission/entity"

type PermissionResponse struct {
	// TODO: add response fields
	Name        string `json:"name"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

func NewPermissionResponse(permission *entity.Permission) *PermissionResponse {
	return &PermissionResponse{
		Name:        permission.Name,
		Group:       permission.Group,
		Description: permission.Description,
	}
}
