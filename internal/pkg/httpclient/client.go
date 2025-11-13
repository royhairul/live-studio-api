package httpclient

import (
	"io"
	"net/http"
)

type Client interface {
	NewRequest(method, rawURL string, query map[string]string, body io.Reader, headers map[string]string) (*http.Request, error)
	DoRequest(req *http.Request, out interface{}) error
}
