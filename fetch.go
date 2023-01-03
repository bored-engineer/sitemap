package sitemap

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"sync"

	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/sync/errgroup"
)

// response captures either a <urlset> or a <sitemapindex>
type response struct {
	XMLName xml.Name
	// Container for the data needed to describe a sitemap.
	Sitemaps []Sitemap `xml:"sitemap"`
	// Container for the data needed to describe a document to crawl.
	URLs []URL `xml:"url"`
}

// fetchResponse performs a single HTTP GET request for the provided URL, parsing the result as a Response struct
func fetchResponse(ctx context.Context, client *http.Client, sitemap string) (*response, error) {
	// If no client was provided, use a reasonable default
	if client == nil {
		client = cleanhttp.DefaultClient()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sitemap, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext failed: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("(*http.Client).Do failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("expected 200, got %d: %s", resp.StatusCode, string(body))
	}
	d := xml.NewDecoder(resp.Body)
	var elem response
	if err := d.Decode(&elem); err != nil {
		return nil, fmt.Errorf("xml.Unmarshal failed: %w", err)
	}
	return &elem, nil
}

// Fetch retrieves the URLs in a given sitemap using reasonable defaults
func Fetch(ctx context.Context, sitemap string, opts ...Option) (urls []URL, err error) {
	// Apply the provided functional options
	var o options
	for _, f := range opts {
		f(&o)
	}
	// Set reasonable defaults if not specified
	if o.client == nil {
		o.client = cleanhttp.DefaultClient()
	}
	if o.parallelism == 0 {
		o.parallelism = runtime.GOMAXPROCS(0)
	}
	if o.processor == nil {
		var mu sync.Mutex
		o.processor = func(_ context.Context, _ *Sitemap, chunk []URL) error {
			mu.Lock()
			urls = append(urls, chunk...)
			mu.Unlock()
			return nil
		}
	}
	// Fetch the root response object first
	resp, err := fetchResponse(ctx, o.client, sitemap)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %q: %w", sitemap, err)
	}
	switch resp.XMLName.Local {
	case "urlset":
		// If it was a <urlset>, we can invoke the processor and exit early
		if err := o.processor(ctx, nil, resp.URLs); err != nil {
			return nil, err
		}
		return
	case "sitemapindex":
		// continue
	default:
		// unexpected root
		return nil, fmt.Errorf("expected <urlset> or <sitemapindex>, got <%s>", resp.XMLName.Local)
	}
	// Must be a <sitemapindex>, parse the provided root URL to extract the hostname for later matching
	rootURL, err := url.Parse(sitemap)
	if err != nil {
		return nil, fmt.Errorf("url.Parse failed: %w", err)
	}
	// If there was a filter func, execute it
	sitemaps := resp.Sitemaps
	if o.filter != nil {
		sitemaps = o.filter(sitemaps)
	}
	// Use an errgroup to bail if an error occurs, and to limit parallelism
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(o.parallelism)
	for _, sitemap := range sitemaps {
		if gCtx.Err() != nil {
			break
		}
		// don't use the loop variables within a go-routine
		sitemap := sitemap
		g.Go(func() error {
			// Check that the origin matches before fetching
			if sitemapURL, err := url.Parse(sitemap.Location); err != nil {
				return fmt.Errorf("url.Parse failed: %w", err)
			} else if sitemapURL.Host != rootURL.Host || sitemapURL.Scheme != rootURL.Scheme {
				return fmt.Errorf("refusing to fetch %q as it is a different origin", sitemap.Location)
			}
			resp, err := fetchResponse(gCtx, o.client, sitemap.Location)
			if err != nil {
				return fmt.Errorf("failed to fetch %q: %w", sitemap.Location, err)
			}
			if resp.XMLName.Local != "urlset" {
				return fmt.Errorf("expected <urlset>, got <%s>", resp.XMLName.Local)
			}
			if err := o.processor(ctx, &sitemap, resp.URLs); err != nil {
				return err
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return
}
