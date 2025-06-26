package shopee

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client struct untuk konfigurasi
type Client struct {
	BaseURL string
	Client  *http.Client
	Headers http.Header
}

// NewClient buat instance baru
func NewClient(baseURL, userCookie string) *Client {
	headers := http.Header{}
	headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	headers.Set("Content-Type", "application/json")
	headers.Set("Cookie", userCookie)

	return &Client{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 15 * time.Second, // optional: timeout 15 detik
		},
		Headers: headers,
	}
}

// RequestOptions untuk flexibility
type RequestOptions struct {
	Method      string
	Endpoint    string
	QueryParams map[string]string
	Body        interface{}
	Headers     map[string]string
}

// DoRequest eksekusi request generik
func (c *Client) DoRequest(opts RequestOptions) ([]byte, error) {
	// Buat URL dengan query
	u, err := url.Parse(fmt.Sprintf("%s/%s", c.BaseURL, opts.Endpoint))
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}
	q := u.Query()
	for key, value := range opts.QueryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	// Marshal body jika ada
	var body io.Reader
	if opts.Body != nil {
		jsonData, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	// Buat request
	req, err := http.NewRequest(opts.Method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header = c.Headers.Clone() // gunakan salinan header default

	// Tambahkan header khusus jika ada
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	// Kirim request
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-2xx status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}
