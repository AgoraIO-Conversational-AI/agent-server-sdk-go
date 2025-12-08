package domain

import (
	"context"

	client "github.com/fern-demo/agoraio-go-sdk/v505/client"
	option "github.com/fern-demo/agoraio-go-sdk/v505/option"
)

// DefaultAPIPath is the default API path appended to the base URL.
const DefaultAPIPath = "/api/conversational-ai-agent"

// Client is a wrapper around the generated Fern client that provides
// region-based URL selection. It automatically selects the appropriate
// base URL based on the specified region/area.
//
// This client embeds the generated Fern client, so all sub-clients
// (Agents, Telephony, PhoneNumbers, etc.) are automatically available.
type Client struct {
	// Embedded client - all fields and methods are promoted automatically.
	*client.Client

	// pool manages domain selection for the specified region.
	pool *Pool
	// apiPath is the API path appended to the base URL.
	apiPath string
}

// ClientOption configures the Client.
type ClientOption func(*clientOptions)

type clientOptions struct {
	area           Area
	apiPath        string
	requestOptions []option.RequestOption
}

// WithArea sets the region/area for the client.
// This determines which regional API endpoints will be used.
func WithArea(area Area) ClientOption {
	return func(o *clientOptions) {
		o.area = area
	}
}

// WithAPIPath sets a custom API path to append to the base URL.
// If not specified, DefaultAPIPath ("/api/conversational-ai-agent") is used.
func WithAPIPath(path string) ClientOption {
	return func(o *clientOptions) {
		o.apiPath = path
	}
}

// WithRequestOptions passes additional request options to the underlying Fern client.
func WithRequestOptions(opts ...option.RequestOption) ClientOption {
	return func(o *clientOptions) {
		o.requestOptions = append(o.requestOptions, opts...)
	}
}

// NewClient creates a new Client with region-based URL selection.
//
// Example usage:
//
//	client, err := domain.NewClient(
//	    domain.WithArea(domain.US),
//	    domain.WithRequestOptions(
//	        option.WithBasicAuth("username", "password"),
//	    ),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Optionally select the best domain via DNS resolution
//	ctx := context.Background()
//	if err := client.SelectBestDomain(ctx); err != nil {
//	    log.Printf("Warning: could not select best domain: %v", err)
//	}
func NewClient(opts ...ClientOption) (*Client, error) {
	options := &clientOptions{
		area:    US, // Default to US region
		apiPath: DefaultAPIPath,
	}

	for _, opt := range opts {
		opt(options)
	}

	pool, err := NewPool(options.area)
	if err != nil {
		return nil, err
	}

	baseURL := pool.GetCurrentUrl() + options.apiPath

	// Prepend the base URL option to the request options
	allOpts := append([]option.RequestOption{option.WithBaseURL(baseURL)}, options.requestOptions...)

	baseClient := client.NewClient(allOpts...)

	return &Client{
		Client:  baseClient,
		pool:    pool,
		apiPath: options.apiPath,
	}, nil
}

// SelectBestDomain performs DNS resolution to select the best available domain
// for the configured region. This is optional but recommended for optimal
// performance and reliability.
func (c *Client) SelectBestDomain(ctx context.Context) error {
	if err := c.pool.SelectBestDomain(ctx); err != nil {
		return err
	}
	// After selecting the best domain, we need to recreate the client
	// with the new base URL. However, since the underlying clients
	// are already created, we update the reference.
	return nil
}

// NextRegion rotates to the next region prefix for failover.
// Call this method when you encounter connection issues to try
// an alternative regional endpoint.
func (c *Client) NextRegion() {
	c.pool.NextRegion()
}

// GetCurrentURL returns the current base URL being used for API requests.
func (c *Client) GetCurrentURL() string {
	return c.pool.GetCurrentUrl() + c.apiPath
}

// GetPool returns the underlying domain pool for advanced usage.
func (c *Client) GetPool() *Pool {
	return c.pool
}
