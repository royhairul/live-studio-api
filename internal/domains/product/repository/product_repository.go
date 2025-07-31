package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/product/entity"
)

type ProductRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Product, error)
	FindByID(id string) (*entity.Product, error)
	Create(data *entity.Product) (*entity.Product, error)
	Update(data *entity.Product) (*entity.Product, error)
	Delete(id string) error
}