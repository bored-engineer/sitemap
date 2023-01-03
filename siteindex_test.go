package sitemap

import (
	"strings"
	"testing"
	"time"

	datetime "github.com/bored-engineer/w3c-datetime"
	"github.com/stretchr/testify/assert"
)

const sampleSitemapIndex = `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<sitemap>
		<loc>http://www.example.com/sitemap1.xml.gz</loc>
		<lastmod>2004-10-01T18:23:17+00:00</lastmod>
	</sitemap>
	<sitemap>
		<loc>http://www.example.com/sitemap2.xml.gz</loc>
		<lastmod>2005-01-01</lastmod>
	</sitemap>
</sitemapindex>`

func TestSitemap(t *testing.T) {
	var s SitemapIndex
	n, err := s.ReadFrom(strings.NewReader(sampleSitemapIndex))
	assert.NoError(t, err)
	assert.Len(t, sampleSitemapIndex, int(n))
	assert.Equal(t, SitemapIndex{
		Sitemaps: []Sitemap{
			{
				Location: "http://www.example.com/sitemap1.xml.gz",
				LastModification: datetime.NewWithPrecision(
					time.Date(2004, 10, 1, 18, 23, 17, 0, time.FixedZone("", 0)),
					datetime.PrecisionSecond,
				),
			},
			{
				Location: "http://www.example.com/sitemap2.xml.gz",
				LastModification: datetime.NewWithPrecision(
					time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC),
					datetime.PrecisionDay,
				),
			},
		},
	}, s)
}
