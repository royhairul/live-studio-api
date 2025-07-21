package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/role/entity"
	"gorm.io/gorm"
)

type RoleRepositoryImpl struct {
	// TODO: add database instance
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &RoleRepositoryImpl{DB: db}
}

// Create implements RoleRepository.
func (r *RoleRepositoryImpl) Create(role *entity.Role) (*entity.Role, error) {
	if err := r.DB.Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

// Delete implements RoleRepository.
func (r *RoleRepositoryImpl) Delete(id string) error {
	if err := r.DB.Delete(entity.Role{}, "id  = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// FindAll implements RoleRepository.
func (r *RoleRepositoryImpl) FindAll() ([]*entity.Role, error) {
	var roles []*entity.Role

	if err := r.DB.Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

// FindByID implements RoleRepository.
func (r *RoleRepositoryImpl) FindByID(id string) (*entity.Role, error) {
	var role entity.Role

	if err := r.DB.Where("id = ?").First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// Update implements RoleRepository.
func (r *RoleRepositoryImpl) Update(role *entity.Role) (*entity.Role, error) {
	if err := r.DB.Model(entity.Role{}).Where("id = ?", role.ID).Updates(&role).Error; err != nil {
		return nil, err
	}

	if err := r.DB.First(role).Error; err != nil {
		return nil, err
	}

	return role, nil
}
