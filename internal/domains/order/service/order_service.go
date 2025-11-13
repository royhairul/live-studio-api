package service

import "github.com/royhairul/live-studio-api/internal/domains/order/params"

type OrderService interface {
	FindAll() ([]*params.OrderResponse, error)
	FindByID(id string) (*params.OrderResponse, error)
	Create(req params.CreateOrderRequest) (*params.OrderResponse, error)
	Update(id string, req params.UpdateOrderRequest) (*params.OrderResponse, error)
	Delete(id string) error
}