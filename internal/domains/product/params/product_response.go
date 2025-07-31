package params

import "github.com/royhairul/live-studio-api/internal/domains/product/entity"

type ProductResponse struct {
	ID       string `json:"id"`
	UniqueID string `json:"unique_id"`
	Name     string `json:"name"`
	ShopID   string `json:"shop_id"`
	ShopName string `json:"shop_name"`
	Platform string `json:"platform"`
	Link     string `json:"link"`
}

func NewProductResponse(product *entity.Product) *ProductResponse {
	return &ProductResponse{
		ID:       product.ID,
		UniqueID: product.UniqueID,
		Name:     product.Name,
		ShopID:   product.ShopID,
		ShopName: product.ShopName,
		Platform: product.Platform,
		Link:     product.Link,
	}
}
