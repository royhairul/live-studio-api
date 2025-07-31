package repository

import (
	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/domains/product/entity"
)

type ProductRepositoryImpl struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &ProductRepositoryImpl{DB: db}
}

// Create implements ProductRepository.
func (r *ProductRepositoryImpl) Create(data *entity.Product) (*entity.Product, error) {
	if err := r.DB.Create(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindAll implements ProductRepository.
func (r *ProductRepositoryImpl) FindAll() ([]*entity.Product, error) {
	var items []*entity.Product
	if err := r.DB.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID implements ProductRepository.
func (r *ProductRepositoryImpl) FindByID(id string) (*entity.Product, error) {
	var item entity.Product
	if err := r.DB.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Update implements ProductRepository.
func (r *ProductRepositoryImpl) Update(data *entity.Product) (*entity.Product, error) {
	if err := r.DB.Model(&entity.Product{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Delete implements ProductRepository.
func (r *ProductRepositoryImpl) Delete(id string) error {
	if err := r.DB.Delete(&entity.Product{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
