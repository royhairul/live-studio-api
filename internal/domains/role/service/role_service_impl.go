package service

import (
	permissionentity "github.com/royhairul/live-studio-api/internal/domains/permission/entity"
	permissionservice "github.com/royhairul/live-studio-api/internal/domains/permission/service"

	"github.com/royhairul/live-studio-api/internal/domains/role/entity"
	"github.com/royhairul/live-studio-api/internal/domains/role/params"
	"github.com/royhairul/live-studio-api/internal/domains/role/repository"
	"gorm.io/gorm"
)

type RoleServiceImpl struct {
	repository    repository.RoleRepository
	permissionSvc permissionservice.PermissionService
}

func NewRoleService(
	repository repository.RoleRepository,
	permissionSvc permissionservice.PermissionService,
) RoleService {
	return &RoleServiceImpl{repository, permissionSvc}
}

// Create implements RoleService.
func (s *RoleServiceImpl) Create(req params.CreateRoleRequest) (*params.RoleResponse, error) {
	role := &entity.Role{
		Name: req.Name,
	}

	for _, id := range req.Permissions {
		role.Permissions = append(role.Permissions, permissionentity.Permission{Model: gorm.Model{ID: id}})
	}

	created, err := s.repository.Create(role)
	if err != nil {
		return nil, err
	}

	result := params.NewRoleResponse(created)
	return result, nil
}

// Update implements RoleService.
func (s *RoleServiceImpl) Update(id string, req params.UpdateRoleRequest) (*params.RoleResponse, error) {
	role, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update nama role
	role.Name = req.Name

	// Reset permissions agar tidak duplikat
	role.Permissions = []permissionentity.Permission{}
	for _, pid := range req.Permissions {
		role.Permissions = append(role.Permissions, permissionentity.Permission{
			Model: gorm.Model{ID: pid},
		})
	}

	updated, err := s.repository.Update(role)
	if err != nil {
		return nil, err
	}

	return params.NewRoleResponse(updated), nil
}

// FindAll implements RoleService.
func (s *RoleServiceImpl) FindAll() ([]*params.RoleResponse, error) {
	roles, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var result []*params.RoleResponse
	for _, role := range roles {
		if role.Name == "superadmin" {
			continue
		}

		result = append(result, params.NewRoleResponse(role))
	}

	return result, nil
}

// FindByID implements RoleService.
func (s *RoleServiceImpl) FindByID(id string) (*params.RoleResponse, error) {
	role, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	result := params.NewRoleResponse(role)
	return result, nil
}

// Delete implements RoleService.
func (s *RoleServiceImpl) Delete(id string) error {
	if err := s.repository.Delete(id); err != nil {
		return err
	}
	return nil
}
