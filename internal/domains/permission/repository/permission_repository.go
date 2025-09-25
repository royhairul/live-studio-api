package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/permission/entity"
)

type PermissionRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Permission, error)
	FindByID(id string) (*entity.Permission, error)
	Create(data *entity.Permission) (*entity.Permission, error)
	Update(data *entity.Permission) (*entity.Permission, error)
	Delete(id string) error
}
