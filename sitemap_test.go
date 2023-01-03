package sitemap

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	u, err := Parse(strings.NewReader(sampleSitemap))
	assert.NoError(t, err)
	assert.IsType(t, URLSet{}, u)
	assert.Len(t, u.(URLSet).URLs, 2)
	s, err := Parse(strings.NewReader(sampleSitemapIndex))
	assert.NoError(t, err)
	assert.IsType(t, SitemapIndex{}, s)
	assert.Len(t, s.(SitemapIndex).Sitemaps, 2)
}
