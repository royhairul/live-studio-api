package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/royhairul/live-studio-api/internal/domains/role/entity"
)

type RoleRepositoryImpl struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &RoleRepositoryImpl{DB: db}
}

// Create implements RoleRepository.
func (r *RoleRepositoryImpl) Create(data *entity.Role) (*entity.Role, error) {
	if err := r.DB.Create(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.Preload(clause.Associations).First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindAll implements RoleRepository.
func (r *RoleRepositoryImpl) FindAll() ([]*entity.Role, error) {
	var items []*entity.Role
	if err := r.DB.Preload(clause.Associations).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID implements RoleRepository.
func (r *RoleRepositoryImpl) FindByID(id string) (*entity.Role, error) {
	var item entity.Role
	if err := r.DB.Preload(clause.Associations).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Update implements RoleRepository.
func (r *RoleRepositoryImpl) Update(data *entity.Role) (*entity.Role, error) {
	// Update kolom role utama
	if err := r.DB.Model(&entity.Role{}).
		Where("id = ?", data.ID).
		Updates(map[string]interface{}{
			"name": data.Name,
		}).Error; err != nil {
		return nil, err
	}

	if err := r.DB.Model(data).
		Association("Permissions").
		Replace(data.Permissions); err != nil {
		return nil, err
	}

	// Reload role lengkap setelah update
	if err := r.DB.Preload(clause.Associations).First(data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// Delete implements RoleRepository.
func (r *RoleRepositoryImpl) Delete(id string) error {
	if err := r.DB.Delete(&entity.Role{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
