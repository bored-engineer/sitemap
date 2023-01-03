package sitemap

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/sync/errgroup"
)

// Fetch performs a single HTTP GET request for the provided URL, parsing the result as a Sitemap or URLSet struct
func Fetch(ctx context.Context, client *http.Client, url string) (any, *http.Response, error) {
	// If no client was provided, use a reasonable default
	if client == nil {
		client = cleanhttp.DefaultClient()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("http.NewRequestWithContext failed: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, resp, fmt.Errorf("(*http.Client).Do failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, resp, fmt.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
	}
	result, err := Parse(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("Parse failed: %w", err)
	}
	return result, resp, nil
}

// FetchParallel fetches the provided URL, if it is already a <urlset> it is returned a single element slice and the returned <sitemapindex> is nil
// If the response was a <sitemap>, each child sitemap is fetched at the provided parallelism and returned in the matching order as the returned <sitemapindex>
func FetchParallel(ctx context.Context, client *http.Client, url string, parallelism int) ([]URLSet, *SitemapIndex, error) {
	// If no client was provided, use a reasonable default
	if client == nil {
		client = cleanhttp.DefaultClient()
	}
	// Fetch the root sitemap
	result, _, err := Fetch(ctx, client, url)
	if err != nil {
		return nil, nil, fmt.Errorf("Fetch for %q failed: %w", url, err)
	}
	// If result was a <urlset>, we're done, bail early
	if urlset, ok := result.(URLSet); ok {
		return []URLSet{urlset}, nil, nil
	}
	// If it was anything other than <sitemapindex>, something has gone wrong
	index, ok := result.(SitemapIndex)
	if !ok {
		return nil, nil, fmt.Errorf("expected SitemapIndex, got %T for %q", result, url)
	}
	// pre-allocate URLs slice so we can (safely) write to it in parallel
	urls := make([]URLSet, len(index.Sitemaps))
	// Use an errgroup to bail if an error occurs, and to limit parallelism
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(parallelism)
	for idx, sitemap := range index.Sitemaps {
		if gCtx.Err() != nil {
			break
		}
		// don't use the loop variables within a go-routine
		idx, url := idx, sitemap.Location
		g.Go(func() error {
			result, _, err := Fetch(gCtx, client, url)
			if err != nil {
				return fmt.Errorf("Fetch for %q failed: %w", url, err)
			}
			// If it was anything other than <urlset>, something has gone wrong
			urlset, ok := result.(URLSet)
			if !ok {
				return fmt.Errorf("expected URLSet, got %T for %q", result, url)
			}
			urls[idx] = urlset
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, nil, err
	}
	return urls, &index, nil
}
