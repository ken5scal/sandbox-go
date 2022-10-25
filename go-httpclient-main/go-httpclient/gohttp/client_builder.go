package gohttp

import (
	"net/http"
	"time"
)

type clientBuilder struct {
	headers            http.Header
	maxIdleConnections int
	disableTimeout     bool
	connectionTimeout  time.Duration
	responseTimeout    time.Duration
	baseUrl            string
	client             *http.Client
	userAgent          string
}

type ClientBuilder interface {
	// Configure the client
	SetHeaders(headers http.Header) ClientBuilder
	SetConnectionTimeout(timeout time.Duration) ClientBuilder
	SetResponseTimeout(timeout time.Duration) ClientBuilder
	SetMaxIdleConnections(maxIdleConnections int) ClientBuilder
	DisableTimeout(disable bool) ClientBuilder
	SetHttpClient(client *http.Client) ClientBuilder
	SetUserAgent(userAgent string) ClientBuilder

	Build() Client
}

func NewBuilder() ClientBuilder {
	// dialer := net.Dialer{Timeout: 1 * time.Second}
	// client := http.Client{
	// 	Transport: &http.Transport{
	// 		MaxIdleConnsPerHost:   5,
	// 		DialContext:           dialer.DialContext,
	// 		ResponseHeaderTimeout: 5 * time.Second,
	// 	},
	// }
	// return &httpClient{client: &client}
	builder := &clientBuilder{}
	return builder
}

func (c *clientBuilder) Build() Client {
	client := httpClient{
		builder: c,
	}
	return &client
}

func (c *clientBuilder) SetHeaders(headers http.Header) ClientBuilder {
	c.headers = headers
	return c
}

func (c *clientBuilder) SetConnectionTimeout(timeout time.Duration) ClientBuilder {
	c.connectionTimeout = timeout
	return c
}

func (c *clientBuilder) SetResponseTimeout(timeout time.Duration) ClientBuilder {
	c.responseTimeout = timeout
	return c
}

func (c *clientBuilder) SetMaxIdleConnections(maxIdleConnections int) ClientBuilder {
	c.maxIdleConnections = maxIdleConnections
	return c
}

func (c *clientBuilder) DisableTimeout(disable bool) ClientBuilder {
	c.disableTimeout = disable
	return c
}

func (c *clientBuilder) SetHttpClient(client *http.Client) ClientBuilder {
	c.client = client
	return c
}

func (c *clientBuilder) SetUserAgent(userAgent string) ClientBuilder {
	c.userAgent = userAgent
	return c
}
