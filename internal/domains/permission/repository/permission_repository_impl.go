package repository

import (
	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/domains/permission/entity"
)

type PermissionRepositoryImpl struct {
	DB *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &PermissionRepositoryImpl{DB: db}
}

// Create implements PermissionRepository.
func (r *PermissionRepositoryImpl) Create(data *entity.Permission) (*entity.Permission, error) {
	if err := r.DB.Create(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindAll implements PermissionRepository.
func (r *PermissionRepositoryImpl) FindAll() ([]*entity.Permission, error) {
	var items []*entity.Permission
	if err := r.DB.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID implements PermissionRepository.
func (r *PermissionRepositoryImpl) FindByID(id string) (*entity.Permission, error) {
	var item entity.Permission
	if err := r.DB.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Update implements PermissionRepository.
func (r *PermissionRepositoryImpl) Update(data *entity.Permission) (*entity.Permission, error) {
	if err := r.DB.Model(&entity.Permission{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Delete implements PermissionRepository.
func (r *PermissionRepositoryImpl) Delete(id string) error {
	if err := r.DB.Delete(&entity.Permission{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// FindByGroup implements PermissionRepository.
func (r *PermissionRepositoryImpl) FindByGroup(group string) ([]*entity.Permission, error) {
	var item []*entity.Permission
	if err := r.DB.Where("group = ?", group).Find(&item).Error; err != nil {
		return nil, err
	}
	return item, nil
}
