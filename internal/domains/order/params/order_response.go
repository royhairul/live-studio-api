package params

import (
	"github.com/royhairul/live-studio-api/internal/domains/order/entity"
	productparams "github.com/royhairul/live-studio-api/internal/domains/product/params"
)

type OrderResponse struct {
	SerialNumber  string                          `json:"serial_number"`
	Status        string                          `json:"status"`
	CompleteTime  int64                           `json:"complete_time"`
	TransactionID int64                           `json:"transaction_id"`
	Products      []productparams.ProductResponse `json:"products"`
}

func NewOrderResponse(order *entity.Order) *OrderResponse {
	var products []productparams.ProductResponse
	for _, product := range order.Products {
		products = append(products, *productparams.NewProductResponse(&product))
	}

	return &OrderResponse{
		SerialNumber:  order.SerialNumber,
		Status:        order.Status,
		CompleteTime:  order.CompleteTime,
		TransactionID: order.TransactionID,
		Products:      products,
	}
}
