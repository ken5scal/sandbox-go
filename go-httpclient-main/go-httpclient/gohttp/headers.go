package gohttp

import (
	"net/http"

	"github.com/ken5scal/go-httpclient/gomime"
)

func getHeaders(headers ...http.Header) http.Header {
	requestHeaders := http.Header{}
	if len(headers) > 0 {
		requestHeaders = headers[0]
	}

	return requestHeaders
}

func (c *httpClient) getRequestHeaders(requestHeaders http.Header) http.Header {
	result := make(http.Header)

	// add common headers to the request
	for header, value := range c.builder.headers {
		if len(value) > 0 {
			result.Set(header, value[0])
		}
	}

	// add custom headers to the requests:
	for header, value := range requestHeaders {
		if len(value) > 0 {
			result.Set(header, value[0])
		}
	}

	if c.builder.userAgent != "" {
		if result.Get(gomime.HeaderUserAgent) != "" {
			return result
		}
		result.Set(gomime.HeaderUserAgent, c.builder.userAgent)
	}
	return result
}
