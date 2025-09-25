package service

import "github.com/royhairul/live-studio-api/internal/domains/permission/params"

type PermissionService interface {
	FindAll() ([]*params.PermissionResponse, error)
	FindByID(id string) (*params.PermissionResponse, error)
	Create(req params.CreatePermissionRequest) (*params.PermissionResponse, error)
	Update(id string, req params.UpdatePermissionRequest) (*params.PermissionResponse, error)
	Delete(id string) error

	FindAllGrouped() ([]*params.PermissionGroupedResponse, error)
}
