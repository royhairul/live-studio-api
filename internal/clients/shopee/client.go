package shopee

import (
	"io"
	"net/http"

	"github.com/royhairul/live-studio-api/internal/pkg/httpclient"
)

type ShopeeClient struct {
	BaseURL string
	Client  httpclient.Client
}

func NewShopeeClient(baseURL string, client httpclient.Client) *ShopeeClient {
	return &ShopeeClient{
		BaseURL: baseURL,
		Client:  client,
	}
}

func (c *ShopeeClient) NewShopeeRequest(method, endpoint string, query map[string]string, body io.Reader, cookie string) (*http.Request, error) {
	headers := map[string]string{
		"User-Agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
		"Content-Type": "application/json",
	}

	if cookie != "" {
		headers["Cookie"] = cookie
	}

	return c.Client.NewRequest(method, c.BaseURL+endpoint, query, body, headers)
}

func (c *ShopeeClient) DoShopeeRequest(req *http.Request, out interface{}) error {
	return c.Client.DoRequest(req, out)
}
