package sitemap

import (
	"context"
	"net/http"
)

// Callback processes each <urlset> as they are identified
type ProcessFunc func(
	ctx context.Context,
	sitemap *Sitemap,
	urls []URL,
) error

// options is an internal
type options struct {
	client      *http.Client
	parallelism int
	processor   ProcessFunc
}

// Option changes the behavior of the Fetch function
type Option func(*options)

// WithHTTPClient replaces the default http.Client (cleanhttp.DefaultClient) used by Fetch
func WithHTTPClient(client *http.Client) Option {
	return func(opts *options) {
		opts.client = client
	}
}

// WithParallelism adjusts the maximum parallel fetches
func WithParallelism(limit int) Option {
	return func(opts *options) {
		opts.parallelism = limit
	}
}

// WithProcessor provides a custom function for processing results
func WithProcessor(f ProcessFunc) Option {
	return func(opts *options) {
		opts.processor = f
	}
}
