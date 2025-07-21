package repository

import "github.com/royhairul/live-studio-api/internal/domains/role/entity"

type RoleRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Role, error)
	FindByID(id string) (*entity.Role, error)
	Create(role *entity.Role) (*entity.Role, error)
	Update(role *entity.Role) (*entity.Role, error)
	Delete(id string) error
}
