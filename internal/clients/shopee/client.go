package shopee

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"
)

type ShopeeClient struct {
	BaseURL string
	HTTP    *http.Client
}

func NewShopeeClient(baseURL string, timeout time.Duration) *ShopeeClient {
	return &ShopeeClient{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *ShopeeClient) NewRequest(method, endpoint string, query map[string]string, body io.Reader, cookie string) (*http.Request, error) {
	fullURL := c.BaseURL + path.Join("/", endpoint)

	// Query Handler
	if query != nil {
		q := url.Values{}
		for k, v := range query {
			q.Set(k, v)
		}

		fullURL += "?" + q.Encode()
	}

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed create request: %w", err)
	}

	// Set Header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", "application/json")

	// If Cookies true
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	return req, nil
}

func (c *ShopeeClient) DoRequest(req *http.Request, out interface{}) error {
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if out != nil {
		if err := json.Unmarshal(bodyBytes, out); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	log.Printf("Response Body: %s\n", string(bodyBytes))

	return nil
}
