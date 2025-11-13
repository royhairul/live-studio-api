package service

import "github.com/royhairul/live-studio-api/internal/domains/product/params"

type ProductService interface {
	FindAll() ([]*params.ProductResponse, error)
	FindByID(id string) (*params.ProductResponse, error)
	Create(req params.CreateProductRequest) (*params.ProductResponse, error)
	Update(id string, req params.UpdateProductRequest) (*params.ProductResponse, error)
	Delete(id string) error
}