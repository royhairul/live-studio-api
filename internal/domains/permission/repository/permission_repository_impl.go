package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/permission/entity"
	"gorm.io/gorm"
)

type PermissionRepositoryImpl struct {
	// TODO: add database instance
	DB *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &PermissionRepositoryImpl{DB: db}
}

// Create implements PermissionRepository.
func (p *PermissionRepositoryImpl) Create(permission *entity.Permission) (*entity.Permission, error) {
	if err := p.DB.Create(permission).Error; err != nil {
		return nil, err
	}

	return permission, nil
}

// Delete implements PermissionRepository.
func (p *PermissionRepositoryImpl) Delete(id string) error {
	if err := p.DB.Delete(entity.Permission{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

// FindAll implements PermissionRepository.
func (p *PermissionRepositoryImpl) FindAll() ([]*entity.Permission, error) {
	var permissions []*entity.Permission

	if err := p.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// FindByID implements PermissionRepository.
func (p *PermissionRepositoryImpl) FindByID(id string) (*entity.Permission, error) {
	var permission entity.Permission

	if err := p.DB.Where("id = ?", id).First(&permission).Error; err != nil {
		return nil, err
	}

	return &permission, nil
}

// Update implements PermissionRepository.
func (p *PermissionRepositoryImpl) Update(permission *entity.Permission) (*entity.Permission, error) {
	if err := p.DB.Model(entity.Permission{}).Where("id = ?", permission.ID).Updates(permission).Error; err != nil {
		return nil, err
	}

	return permission, nil
}
