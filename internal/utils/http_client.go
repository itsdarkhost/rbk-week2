package utils

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

var defaultHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
}

type Client struct {
	baseURL      string
	http         *http.Client
	defaultQuery map[string]string
}

func NewHTTPClient(timeout time.Duration) *http.Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &http.Client{
		Timeout: timeout,
	}
}

func NewClient(baseURL string, httpClient *http.Client, defaultQuery map[string]string) *Client {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}

	if defaultQuery == nil {
		defaultQuery = map[string]string{}
	}

	return &Client{
		baseURL:      baseURL,
		http:         httpClient,
		defaultQuery: defaultQuery,
	}
}

func (c *Client) newRequest(ctx context.Context, method, path string, query map[string]string) (*http.Request, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}

	q := u.Query()

	for k, v := range c.defaultQuery {
		if v != "" {
			q.Set(k, v)
		}
	}

	for k, v := range query {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()

	return http.NewRequestWithContext(ctx, method, u.String(), nil)
}

func (c *Client) Get(ctx context.Context, path string, query map[string]string) (*http.Response, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	return c.http.Do(req)
}
