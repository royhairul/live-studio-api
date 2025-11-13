package params

import "github.com/royhairul/live-studio-api/internal/domains/permission/entity"

type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" gorm:"unique"`
	Group       string `json:"group"`
	Description string `json:"description"`
}

type PermissionGroupedResponse struct {
	Group       string                `json:"group"`
	Permissions []*PermissionResponse `json:"permissions"`
}

func NewPermissionResponse(permission *entity.Permission) *PermissionResponse {
	return &PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Group:       permission.Group,
		Description: permission.Description,
	}
}
