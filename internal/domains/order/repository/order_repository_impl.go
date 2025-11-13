package repository

import (
	"gorm.io/gorm"

	"github.com/royhairul/live-studio-api/internal/domains/order/entity"
)

type OrderRepositoryImpl struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &OrderRepositoryImpl{DB: db}
}

// Create implements OrderRepository.
func (r *OrderRepositoryImpl) Create(data *entity.Order) (*entity.Order, error) {
	if err := r.DB.Create(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// FindAll implements OrderRepository.
func (r *OrderRepositoryImpl) FindAll() ([]*entity.Order, error) {
	var items []*entity.Order
	if err := r.DB.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID implements OrderRepository.
func (r *OrderRepositoryImpl) FindByID(id string) (*entity.Order, error) {
	var item entity.Order
	if err := r.DB.Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Update implements OrderRepository.
func (r *OrderRepositoryImpl) Update(data *entity.Order) (*entity.Order, error) {
	if err := r.DB.Model(&entity.Order{}).Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return nil, err
	}
	if err := r.DB.First(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// Delete implements OrderRepository.
func (r *OrderRepositoryImpl) Delete(id string) error {
	if err := r.DB.Delete(&entity.Order{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
