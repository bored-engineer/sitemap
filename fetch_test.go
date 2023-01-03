package sitemap

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "sitemap") {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, sampleSitemap)
		} else {
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "not found")
		}
	}))
	defer ts.Close()
	ctx := context.TODO()
	result, _, err := Fetch(ctx, nil, ts.URL+"/sitemap")
	assert.NoError(t, err)
	assert.IsType(t, URLSet{}, result)
	result, _, err = Fetch(ctx, nil, ts.URL+"/404")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFetchParallel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		baseURL := "http://" + r.Host + "/"
		if strings.Contains(r.URL.Path, "index") {
			io.WriteString(w, strings.ReplaceAll(sampleSitemapIndex, "http://www.example.com/", baseURL))
		} else {
			io.WriteString(w, strings.ReplaceAll(sampleSitemap, "http://www.example.com/", baseURL))
		}
	}))
	defer ts.Close()
	ctx := context.TODO()
	result, index, err := FetchParallel(ctx, nil, ts.URL+"/urlset", -1)
	assert.NoError(t, err)
	assert.Nil(t, index)
	assert.Len(t, result, 1)
	result, index, err = FetchParallel(ctx, nil, ts.URL+"/index", 2)
	assert.NoError(t, err)
	assert.NotNil(t, index)
	assert.Len(t, result, 2)
	assert.Len(t, index.Sitemaps, 2)
}
