package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/product/entity"
	"github.com/royhairul/live-studio-api/internal/domains/product/params"
	"github.com/royhairul/live-studio-api/internal/domains/product/repository"
)

type ProductServiceImpl struct {
	repository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) ProductService {
	return &ProductServiceImpl{repository}
}

// Create implements ProductService.
func (s *ProductServiceImpl) Create(req params.CreateProductRequest) (*params.ProductResponse, error) {
	product := entity.Product{
		UniqueID: req.UniqueID,
		Name:     req.Name,
		ShopID:   req.ShopID,
		ShopName: req.ShopName,
		Platform: req.Platform,
		Link:     req.Link,
	}

	created, err := s.repository.Create(&product)
	if err != nil {
		return nil, err
	}

	result := params.NewProductResponse(created)
	return result, nil
}

// Update implements ProductService.
func (s *ProductServiceImpl) Update(id string, req params.UpdateProductRequest) (*params.ProductResponse, error) {
	panic("unimplemented")
}

// FindAll implements ProductService.
func (s *ProductServiceImpl) FindAll() ([]*params.ProductResponse, error) {
	panic("unimplemented")
}

// FindByID implements ProductService.
func (s *ProductServiceImpl) FindByID(id string) (*params.ProductResponse, error) {
	panic("unimplemented")
}

// Delete implements ProductService.
func (s *ProductServiceImpl) Delete(id string) error {
	panic("unimplemented")
}
