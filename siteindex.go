package sitemap

import (
	"encoding/xml"
	"io"

	datetime "github.com/bored-engineer/w3c-datetime"
)

// Ensure SitemapIndex implements io.ReaderFrom at compile
var _ io.ReaderFrom = (*SitemapIndex)(nil)

// Container for the data needed to describe a document to crawl.
type Sitemap struct {
	// REQUIRED: The location URI of a document.
	// The URI must conform to RFC 2396 (http://www.ietf.org/rfc/rfc2396.txt).
	Location string `xml:"loc"`
	// OPTIONAL: The date the document was last modified.
	// The date must conform to the W3C DATETIME format (http://www.w3.org/TR/NOTE-datetime).
	// Example: 2005-05-10 Lastmod may also contain a timestamp. Example: 2005-05-10T17:33:30+08:00
	LastModification datetime.Time `xml:"lastmod,omitempty"`
}

// Container for a set of up to 50,000 sitemap URLs. This is the root element of the XML file.
type SitemapIndex struct {
	// Container for the data needed to describe a sitemap.
	Sitemaps []Sitemap `xml:"sitemap"`
}

// ReadFrom implements io.ReaderFrom.
func (s *SitemapIndex) ReadFrom(r io.Reader) (int64, error) {
	de := xml.NewDecoder(r)
	err := de.Decode(s)
	return de.InputOffset(), err
}
