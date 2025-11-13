package repository

import (
	"github.com/royhairul/live-studio-api/internal/domains/order/entity"
)

type OrderRepository interface {
	// TODO: define repository methods
	FindAll() ([]*entity.Order, error)
	FindByID(id string) (*entity.Order, error)
	Create(data *entity.Order) (*entity.Order, error)
	Update(data *entity.Order) (*entity.Order, error)
	Delete(id string) error
}