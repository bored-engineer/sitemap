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

func TestFetch_URLSet(t *testing.T) {
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
	urls, err := Fetch(ctx, ts.URL+"/sitemap")
	assert.NoError(t, err)
	assert.Len(t, urls, 2)
	urls, err = Fetch(ctx, ts.URL+"/404")
	assert.Error(t, err)
	assert.Nil(t, urls)
}

func TestFetch_SitemapIndex(t *testing.T) {
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
	urls, err := Fetch(ctx, ts.URL+"/urlset")
	assert.NoError(t, err)
	assert.Len(t, urls, 2)
	urls, err = Fetch(ctx, ts.URL+"/index")
	assert.NoError(t, err)
	assert.Len(t, urls, 4)
}
