package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/permission/entity"
	"github.com/royhairul/live-studio-api/internal/domains/permission/params"
	"github.com/royhairul/live-studio-api/internal/domains/permission/repository"
)

type PermissionServiceImpl struct {
	// TODO: add repository dependency
	repository repository.PermissionRepository
}

func NewPermissionService(repository repository.PermissionRepository) PermissionService {
	return &PermissionServiceImpl{repository}
}

// Create implements PermissionService.
func (p *PermissionServiceImpl) Create(permissionReq params.CreatePermissionRequest) (*params.PermissionResponse, error) {
	permission := entity.Permission{
		Name:        permissionReq.Name,
		Group:       permissionReq.Group,
		Description: permissionReq.Description,
	}

	created, err := p.repository.Create(&permission)
	if err != nil {
		return nil, err
	}

	result := params.NewPermissionResponse(created)

	return result, nil
}

// Delete implements PermissionService.
func (p *PermissionServiceImpl) Delete(id string) error {
	if err := p.repository.Delete(id); err != nil {
		return err
	}

	return nil
}

// FindAll implements PermissionService.
func (p *PermissionServiceImpl) FindAll() ([]*params.PermissionResponse, error) {
	permissions, err := p.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var results []*params.PermissionResponse
	for _, permission := range permissions {
		results = append(results, params.NewPermissionResponse(permission))
	}
	return results, nil
}

// FindByID implements PermissionService.
func (p *PermissionServiceImpl) FindByID(id string) (*params.PermissionResponse, error) {
	permission, err := p.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	result := params.NewPermissionResponse(permission)
	return result, nil
}

// Update implements PermissionService.
func (p *PermissionServiceImpl) Update(permissionReq params.UpdatePermissionRequest) (*params.PermissionResponse, error) {
	panic("unimplemented")
}
