package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/role/params"
	"github.com/royhairul/live-studio-api/internal/domains/role/repository"
)

type RoleServiceImpl struct {
	// TODO: add repository dependency
	repository repository.RoleRepository
}

func NewRoleService(repository repository.RoleRepository) RoleService {
	return &RoleServiceImpl{repository}
}

// Create implements RoleService.
func (r *RoleServiceImpl) Create(roleReq params.CreateRoleRequest) (*params.RoleResponse, error) {
	panic("unimplemented")
}

// Delete implements RoleService.
func (r *RoleServiceImpl) Delete(id string) error {
	if err := r.repository.Delete(id); err != nil {
		return err
	}

	return nil
}

// FindAll implements RoleService.
func (r *RoleServiceImpl) FindAll() ([]*params.RoleResponse, error) {
	roles, err := r.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var result []*params.RoleResponse
	for _, role := range roles {
		result = append(result, params.NewRoleResponse(role))
	}

	return result, nil
}

// FindByID implements RoleService.
func (r *RoleServiceImpl) FindByID(id string) (*params.RoleResponse, error) {
	role, err := r.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	result := params.NewRoleResponse(role)
	return result, nil
}

// Update implements RoleService.
func (r *RoleServiceImpl) Update(id string, roleReq params.UpdateRoleRequest) (*params.RoleResponse, error) {
	panic("unimplemented")
}
