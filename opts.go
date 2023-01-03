package sitemap

import (
	"context"
	"net/http"
)

// ProcessFunc processes each <urlset> as they are identified
type ProcessFunc func(
	ctx context.Context,
	sitemap *Sitemap,
	urls []URL,
) error

// FilterFunc reduces the list of sitemaps to fetch
type FilterFunc func([]Sitemap) []Sitemap

// options is an internal
type options struct {
	client      *http.Client
	parallelism int
	processor   ProcessFunc
	filter      FilterFunc
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

// WithFilter provides a custom function for filtering sitemaps
func WithFilter(f FilterFunc) Option {
	return func(opts *options) {
		opts.filter = f
	}
}
