package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/permission/entity"
	"github.com/royhairul/live-studio-api/internal/domains/permission/params"
	"github.com/royhairul/live-studio-api/internal/domains/permission/repository"
)

type PermissionServiceImpl struct {
	repository repository.PermissionRepository
}

func NewPermissionService(repository repository.PermissionRepository) PermissionService {
	return &PermissionServiceImpl{repository}
}

// Create implements PermissionService.
func (s *PermissionServiceImpl) Create(req params.CreatePermissionRequest) (*params.PermissionResponse, error) {
	permission := &entity.Permission{
		Name:        req.Name,
		Group:       req.Group,
		Description: req.Description,
	}

	created, err := s.repository.Create(permission)
	if err != nil {
		return nil, err
	}

	result := params.NewPermissionResponse(created)
	return result, nil
}

// Update implements PermissionService.
func (s *PermissionServiceImpl) Update(id string, req params.UpdatePermissionRequest) (*params.PermissionResponse, error) {
	panic("unimplemented")
}

// FindAll implements PermissionService.
func (s *PermissionServiceImpl) FindAll() ([]*params.PermissionResponse, error) {
	permissions, err := s.repository.FindAll()
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
func (s *PermissionServiceImpl) FindByID(id string) (*params.PermissionResponse, error) {
	permission, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}

	result := params.NewPermissionResponse(permission)
	return result, nil
}

// Delete implements PermissionService.
func (s *PermissionServiceImpl) Delete(id string) error {
	panic("unimplemented")
}

// FindByGroup implements PermissionService.
func (s *PermissionServiceImpl) FindByGroup(group string) ([]*params.PermissionResponse, error) {
	permissions, err := s.repository.FindByGroup(group)
	if err != nil {
		return nil, err
	}

	var results []*params.PermissionResponse
	for _, permission := range permissions {
		results = append(results, params.NewPermissionResponse(permission))
	}
	return results, nil
}
