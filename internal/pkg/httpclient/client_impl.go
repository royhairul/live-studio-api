package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ClientImpl struct {
	HTTP *http.Client
}

func NewClient(httpClient *http.Client) Client {
	return &ClientImpl{HTTP: httpClient}
}

// NewRequest implements Client.
func (c *ClientImpl) NewRequest(method string, rawURL string, query map[string]string, body io.Reader, headers map[string]string) (*http.Request, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	q := parsedURL.Query()
	for key, val := range query {
		q.Set(key, val)
	}
	parsedURL.RawQuery = q.Encode()

	req, err := http.NewRequest(method, parsedURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

// DoRequest implements Client.
func (c *ClientImpl) DoRequest(req *http.Request, out interface{}) error {
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if out != nil {
		if err := json.Unmarshal(bodyBytes, out); err != nil {
			return fmt.Errorf("unmarshal response failed: %w", err)
		}
	}

	return nil
}
