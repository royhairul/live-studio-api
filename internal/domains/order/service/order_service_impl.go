package service

import (
	"github.com/royhairul/live-studio-api/internal/domains/order/entity"
	"github.com/royhairul/live-studio-api/internal/domains/order/params"
	"github.com/royhairul/live-studio-api/internal/domains/order/repository"
	productentity "github.com/royhairul/live-studio-api/internal/domains/product/entity"
	productservice "github.com/royhairul/live-studio-api/internal/domains/product/service"
)

type OrderServiceImpl struct {
	repository     repository.OrderRepository
	productService productservice.ProductService
}

func NewOrderService(repository repository.OrderRepository, productService productservice.ProductService) OrderService {
	return &OrderServiceImpl{repository, productService}
}

// Create implements OrderService.
func (s *OrderServiceImpl) Create(req params.CreateOrderRequest) (*params.OrderResponse, error) {
	var products []productentity.Product
	for _, productReq := range req.Products {
		productResp, err := s.productService.Create(productReq)
		if err != nil {
			return nil, err
		}

		// Convert ProductResponse to entity.Product
		productEntity := productentity.Product{
			ID:       productResp.ID,
			Name:     productResp.Name,
			UniqueID: productResp.UniqueID,
			ShopID:   productResp.ShopID,
			ShopName: productResp.ShopName,
			// Add other fields as necessary
		}

		products = append(products, productEntity)
	}

	order := entity.Order{
		SerialNumber:  req.SerialNumber,
		Status:        req.Status,
		CompleteTime:  req.CompleteTime,
		TransactionID: req.TransactionID,
		Products:      products,
	}

	created, err := s.repository.Create(&order)
	if err != nil {
		return nil, err
	}

	result := params.NewOrderResponse(created)
	return result, nil
}

// Update implements OrderService.
func (s *OrderServiceImpl) Update(id string, req params.UpdateOrderRequest) (*params.OrderResponse, error) {
	panic("unimplemented")
}

// FindAll implements OrderService.
func (s *OrderServiceImpl) FindAll() ([]*params.OrderResponse, error) {
	orders, err := s.repository.FindAll()
	if err != nil {
		return nil, err
	}

	var result []*params.OrderResponse
	for _, order := range orders {
		result = append(result, params.NewOrderResponse(order))
	}

	return result, nil
}

// FindByID implements OrderService.
func (s *OrderServiceImpl) FindByID(id string) (*params.OrderResponse, error) {
	panic("unimplemented")
}

// Delete implements OrderService.
func (s *OrderServiceImpl) Delete(id string) error {
	panic("unimplemented")
}
