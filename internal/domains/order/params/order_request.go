package params

import (
	productparams "github.com/royhairul/live-studio-api/internal/domains/product/params"
)

type OrderRequest struct {
	// TODO: add request fields
}

type CreateOrderRequest struct {
	SerialNumber  string                               `json:"serial_number"`
	Status        string                               `json:"status"`
	CompleteTime  int64                                `json:"complete_time"`
	TransactionID int64                                `json:"transaction_id"`
	Products      []productparams.CreateProductRequest `json:"products"`
}

type UpdateOrderRequest struct {
	// TODO: add request fields
}
