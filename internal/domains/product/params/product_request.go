package params

type ProductRequest struct {
	// TODO: add request fields
}

type CreateProductRequest struct {
	UniqueID string `json:"unique_id"`
	Name     string `json:"name"`
	ShopID   string `json:"shop_id"`
	ShopName string `json:"shop_name"`
	Platform string `json:"platform"`
	Link     string `json:"link"`
}

type UpdateProductRequest struct {
	// TODO: add request fields
}
