package sitemap

import (
	"encoding/xml"
	"io"

	datetime "github.com/bored-engineer/w3c-datetime"
)

// Ensure URLSet implements io.ReaderFrom at compile
var _ io.ReaderFrom = (*URLSet)(nil)

// Alternate is a map of language codes to corrosponding URL
// https://developers.google.com/search/docs/specialty/international/localized-versions
type Alternate map[string]string

// UnmarshalXML implements xml.Unmarshaler
func (a *Alternate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var link struct {
		Relationship string `xml:"rel,attr"`
		Language     string `xml:"hreflang,attr"`
		URL          string `xml:"href,attr"`
	}
	if err := d.DecodeElement(&link, &start); err != nil {
		return err
	}
	if *a == nil {
		*a = make(Alternate)
	}
	if link.Relationship == "alternate" {
		(*a)[link.Language] = link.URL
	}
	return nil
}

// OPTIONAL: Indicates how frequently the content at a particular URL is likely to change.
// The value "always" should be used to describe documents that change each time they are accessed.
// The value "never" should be used to describe archived URLs.
// Please note that web crawlers may not necessarily crawl pages marked "always" more often.
// Consider this element as a friendly suggestion and not a command.
type ChangeFrequency string

const (
	ChangeFrequencyAlways  ChangeFrequency = "always"
	ChangeFrequencyHourly  ChangeFrequency = "hourly"
	ChangeFrequencyDaily   ChangeFrequency = "daily"
	ChangeFrequencyWeekly  ChangeFrequency = "weekly"
	ChangeFrequencyMonthly ChangeFrequency = "monthly"
	ChangeFrequencyYearly  ChangeFrequency = "yearly"
	ChangeFrequencyNever   ChangeFrequency = "never"
)

// Container for the data needed to describe a document to crawl.
type URL struct {
	// REQUIRED: The location URI of a document.
	// The URI must conform to RFC 2396 (http://www.ietf.org/rfc/rfc2396.txt).
	Location string `xml:"loc"`
	// OPTIONAL: The date the document was last modified.
	// The date must conform to the W3C DATETIME format (http://www.w3.org/TR/NOTE-datetime).
	// Example: 2005-05-10 Lastmod may also contain a timestamp. Example: 2005-05-10T17:33:30+08:00
	LastModification datetime.Time `xml:"lastmod,omitempty"`
	// OPTIONAL: Indicates how frequently the content at a particular URL is likely to change.
	// The value "always" should be used to describe documents that change each time they are accessed.
	// The value "never" should be used to describe archived URLs.
	// Please note that web crawlers may not necessarily crawl pages marked "always" more often.
	// Consider this element as a friendly suggestion and not a command.
	ChangeFrequency ChangeFrequency `xml:"changefreq,omitempty"`
	// OPTIONAL: The priority of a particular URL relative to other pages on the same site.
	// The value for this element is a number between 0.0 and 1.0 where 0.0 identifies the lowest priority page(s).
	// The default priority of a page is 0.5. Priority is used to select between pages on your site.
	// Setting a priority of 1.0 for all URLs will not help you, as the relative priority of pages on your site is what will be considered.
	Priority float64 `xml:"priority,omitempty"`
	// https://developers.google.com/search/docs/specialty/international/localized-versions
	Alternate Alternate `xml:"link,omitempty"`
}

// Container for a set of up to 50,000 document elements. This is the root element of the XML file.
type URLSet struct {
	// Container for the data needed to describe a document to crawl.
	URLs []URL `xml:"url"`
}

// ReadFrom implements io.ReaderFrom.
func (u *URLSet) ReadFrom(r io.Reader) (int64, error) {
	de := xml.NewDecoder(r)
	err := de.Decode(u)
	return de.InputOffset(), err
}
