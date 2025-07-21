package service

import "github.com/royhairul/live-studio-api/internal/domains/role/params"

type RoleService interface {
	// TODO: define service methods
	FindAll() ([]*params.RoleResponse, error)
	FindByID(id string) (*params.RoleResponse, error)
	Create(roleReq params.CreateRoleRequest) (*params.RoleResponse, error)
	Update(id string, roleReq params.UpdateRoleRequest) (*params.RoleResponse, error)
	Delete(id string) error
}
