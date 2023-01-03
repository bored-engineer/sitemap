package sitemap

import (
	"strings"
	"testing"
	"time"

	datetime "github.com/bored-engineer/w3c-datetime"
	"github.com/stretchr/testify/assert"
)

const sampleSitemap = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">
	<url>
		<loc>http://www.example.com/bar</loc>
		<lastmod>2005-01-01</lastmod>
		<changefreq>monthly</changefreq>
		<priority>0.8</priority>
	</url>
	<url>
		<loc>http://www.example.com/foo</loc>
		<changefreq>weekly</changefreq>
		<xhtml:link rel="alternate" hreflang="de" href="http://www.example.com/foo?hl=de"/>
		<xhtml:link rel="alternate" hreflang="de-ch" href="http://www.example.com/foo?hl=de-ch"/>
		<xhtml:link rel="alternate" hreflang="en" href="http://www.example.com/foo?hl=en"/>
	</url>
</urlset>`

func TestURLSet(t *testing.T) {
	var u URLSet
	n, err := u.ReadFrom(strings.NewReader(sampleSitemap))
	assert.NoError(t, err)
	assert.Len(t, sampleSitemap, int(n))
	assert.Equal(t, URLSet{
		URLs: []URL{
			{
				Location: "http://www.example.com/bar",
				LastModification: datetime.NewWithPrecision(
					time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC),
					datetime.PrecisionDay,
				),
				ChangeFrequency: ChangeFrequencyMonthly,
				Priority:        0.8,
			},
			{
				Location:        "http://www.example.com/foo",
				ChangeFrequency: ChangeFrequencyWeekly,
				Alternate: Alternate{
					"de":    "http://www.example.com/foo?hl=de",
					"de-ch": "http://www.example.com/foo?hl=de-ch",
					"en":    "http://www.example.com/foo?hl=en",
				},
			},
		},
	}, u)
}
