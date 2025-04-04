package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an HTTP client with configurable options
type Client struct {
	HTTPClient     *http.Client
	DefaultHeaders map[string]string
}

// RequestOptions contains options for making HTTP requests
type RequestOptions struct {
	Headers map[string]string
	Query   map[string]string
	Body    interface{}
}

// Response represents an HTTP response
type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
}

// NewClient creates a new HTTP client with default configuration
func NewClient(baseURL string) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		DefaultHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}
}

// Get performs an HTTP GET request
func (c *Client) Get(ctx context.Context, url string, opts *RequestOptions) (*Response, error) {
	return c.Request(ctx, http.MethodGet, url, opts)
}

// Post performs an HTTP POST request
func (c *Client) Post(ctx context.Context, url string, opts *RequestOptions) (*Response, error) {
	return c.Request(ctx, http.MethodPost, url, opts)
}

// Put performs an HTTP PUT request
func (c *Client) Put(ctx context.Context, url string, opts *RequestOptions) (*Response, error) {
	return c.Request(ctx, http.MethodPut, url, opts)
}

// Delete performs an HTTP DELETE request
func (c *Client) Delete(ctx context.Context, url string, opts *RequestOptions) (*Response, error) {
	return c.Request(ctx, http.MethodDelete, url, opts)
}

// Request performs an HTTP request with the given method, url, and options
func (c *Client) Request(ctx context.Context, method, url string, opts *RequestOptions) (*Response, error) {
	if opts == nil {
		opts = &RequestOptions{}
	}

	// Prepare request body
	var reqBody io.Reader
	if opts.Body != nil {
		jsonBody, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add query parameters
	q := req.URL.Query()
	for key, value := range opts.Query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// Add default headers
	for key, value := range c.DefaultHeaders {
		req.Header.Set(key, value)
	}

	// Add request-specific headers
	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
		Body:       body,
	}, nil
}

// DecodeJSON decodes JSON response body into the given target
func (r *Response) DecodeJSON(target interface{}) error {
	return json.Unmarshal(r.Body, target)
}

// IsSuccess returns true if the response status code is in the 2xx range
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}
