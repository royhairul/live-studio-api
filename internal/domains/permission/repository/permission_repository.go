package repository

import "github.com/royhairul/live-studio-api/internal/domains/permission/entity"

type PermissionRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Permission, error)
	FindByID(id string) (*entity.Permission, error)
	Create(permission *entity.Permission) (*entity.Permission, error)
	Update(permission *entity.Permission) (*entity.Permission, error)
	Delete(id string) error
}
