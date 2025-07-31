package service

import "github.com/royhairul/live-studio-api/internal/domains/role/params"

type RoleService interface {
	FindAll() ([]*params.RoleResponse, error)
	FindByID(id string) (*params.RoleResponse, error)
	Create(req params.CreateRoleRequest) (*params.RoleResponse, error)
	Update(id string, req params.UpdateRoleRequest) (*params.RoleResponse, error)
	Delete(id string) error
}