package sitemap

import (
	"encoding/xml"
	"io"
)

// Parse consumes from an io.Reader and returns a populated URLSet or SitemapIndex struct, without consuming the reader twice
func Parse(r io.Reader) (result any, err error) {
	d := xml.NewDecoder(r)
	for {
		// Consume until we reach EOF
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				// If no matching root element was found, unexpected EOF error
				if result == nil {
					return nil, io.ErrUnexpectedEOF
				} else {
					return result, nil
				}
			}
			return nil, err
		}
		elem, ok := t.(xml.StartElement)
		if !ok {
			continue
		}
		// Depending on the local name, result type is determined (or element is skipped)
		switch elem.Name.Local {
		case "urlset":
			var u URLSet
			if err := d.DecodeElement(&u, &elem); err != nil {
				return nil, err
			}
			result = u
		case "sitemapindex":
			var s SitemapIndex
			if err := d.DecodeElement(&s, &elem); err != nil {
				return nil, err
			}
			result = s
		default:
			// Skip over unknown elements
			if err := d.Skip(); err != nil {
				return nil, err
			}
		}
	}
}
