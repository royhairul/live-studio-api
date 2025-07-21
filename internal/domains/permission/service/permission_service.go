package service

import "github.com/royhairul/live-studio-api/internal/domains/permission/params"

type PermissionService interface {
	// TODO: define service methods
	FindAll() ([]*params.PermissionResponse, error)
	FindByID(id string) (*params.PermissionResponse, error)
	Create(permissionReq params.CreatePermissionRequest) (*params.PermissionResponse, error)
	Update(permissionReq params.UpdatePermissionRequest) (*params.PermissionResponse, error)
	Delete(id string) error
}
