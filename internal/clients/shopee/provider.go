package shopee

import "github.com/royhairul/live-studio-api/internal/pkg/httpclient"

func NewShopeeCreatorClient(client httpclient.Client) *ShopeeClient {
	return &ShopeeClient{
		Client:  client,
		BaseURL: "https://creator.shopee.co.id",
	}
}

func NewShopeeSellerClient(client httpclient.Client) *ShopeeClient {
	return &ShopeeClient{
		Client:  client,
		BaseURL: "https://seller.shopee.co.id",
	}
}

func NewShopeeAffiliateClient(client httpclient.Client) *ShopeeClient {
	return &ShopeeClient{
		Client:  client,
		BaseURL: "https://affiliate.shopee.co.id",
	}
}

func NewShopeeDefaultClient(client httpclient.Client) *ShopeeClient {
	return &ShopeeClient{
		Client:  client,
		BaseURL: "https://shopee.co.id",
	}
}
