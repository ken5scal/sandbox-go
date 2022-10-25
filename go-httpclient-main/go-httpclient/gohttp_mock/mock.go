package gohttp_mock

import (
	"fmt"
	"net/http"

	"github.com/ken5scal/go-httpclient/core"
)

// The Mock structure provides a xlean way to configure HTTP mocks based on
// the combination between request method, URL and request body.
type Mock struct {
	Method      string
	Url         string
	RequestBody string

	Error              error
	ResponseBody       string
	ResponseStatusCode int
}

// GetResponse returns a response object based on the mock configuration.
func (m *Mock) GetResponse() (*core.Response, error) {
	if m.Error != nil {
		return nil, m.Error
	}

	response := core.Response{
		Status:     fmt.Sprintf("%d %s", m.ResponseStatusCode, http.StatusText((m.ResponseStatusCode))),
		StatusCode: m.ResponseStatusCode,
		Body:       []byte(m.ResponseBody),
	}
	return &response, nil
}
