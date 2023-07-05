package atom

import (
	"encoding/xml"
	"fmt"
	"time"
)

type Link struct {
	// Required elements
	Href string `xml:"href,attr"`

	// Optional elements
	Rel      string `xml:"rel,attr,omitempty"`
	Type     string `xml:"type,attr,omitempty"`
	HrefLang string `xml:"hreflang,attr,omitempty"`
	Title    string `xml:"title,attr,omitempty"`
	Length   int    `xml:"length,attr,omitempty"`
}

// NewLink creates a new Atom link.
func NewLink(href string) *Link {
	return &Link{Href: href}
}

// SetRel sets the rel attribute of the link.
//
// Available values are:
//
// - alternate: an alternate representation of the entry or feed, for example a permalink to the html version of the entry, or the front page of the weblog.
//
// - enclosure: a related resource which is potentially large in size and might require special handling, for example an audio or video recording.
//
// - related: an document related to the entry or feed.
//
// - self: the feed itself.
//
// - via: the source of the information provided in the entry.
func (l *Link) SetRel(rel string) *Link {
	l.Rel = rel
	return l
}

// SetType indicates the media type of the resource.
func (l *Link) SetType(t string) *Link {
	l.Type = t
	return l
}

// SetHrefLang indicates the language of the referenced resource.
func (l *Link) SetHrefLang(lang string) *Link {
	l.HrefLang = lang
	return l
}

// SetTitle sets human readable information about the link, typically for display purposes.
func (l *Link) SetTitle(title string) *Link {
	l.Title = title
	return l
}

// SetLength sets the length of the resource, in bytes.
func (l *Link) SetLength(length int) *Link {
	l.Length = length
	return l
}

type Person struct {
	// Required elements
	Name string `xml:"name"`

	// Optional elements
	Email string `xml:"email,omitempty"`
	URI   string `xml:"uri,omitempty"`
}

// NewPerson describes a person, corporation, or similar entity in an Atom feed.
func NewPerson(name string) *Person {
	return &Person{Name: name}
}

// SetEmail sets the email address of the person.
func (p *Person) SetEmail(email string) *Person {
	p.Email = email
	return p
}

// SetURI sets a homepage for the person.
func (p *Person) SetURI(uri string) *Person {
	p.URI = uri
	return p
}

type Category struct {
	// Required elements
	Term string `xml:"term,attr"`

	// Optional elements
	Scheme string `xml:"scheme,attr,omitempty"`
	Label  string `xml:"label,attr,omitempty"`
}

// NewCategory creates a new Atom category.
func NewCategory(term string) *Category {
	return &Category{Term: term}
}

// SetScheme sets the the categorization scheme via a URI.
func (c *Category) SetScheme(scheme string) *Category {
	c.Scheme = scheme
	return c
}

// SetLabel sets a human-readable label for display.
func (c *Category) SetLabel(label string) *Category {
	c.Label = label
	return c
}

type Generator struct {
	// Required elements
	Text string `xml:",chardata"`

	// Optional elements
	URI     string `xml:"uri,attr"`
	Version string `xml:"version,attr"`
}

// NewGenerator creates a new Atom generator.
func NewGenerator(text string) *Generator {
	return &Generator{Text: text}
}

// SetURI sets the URI of the generating agent.
func (g *Generator) SetURI(uri string) *Generator {
	g.URI = uri
	return g
}

// SetVersion sets the version of the generating agent.
func (g *Generator) SetVersion(version string) *Generator {
	g.Version = version
	return g
}

type Source struct {
	// Required elements
	ID string `xml:"id"`

	// Optional elements
	Title   string `xml:"title"`
	Updated string `xml:"updated"`
}

// NewSource creates a new Atom source.
func NewSource(id string) *Source {
	return &Source{ID: id}
}

// SetTitle sets the human readable title of the source.
func (s *Source) SetTitle(title string) *Source {
	s.Title = title
	return s
}

// SetUpdated sets the last time the source was modified.
func (s *Source) SetUpdated(updated time.Time) *Source {
	s.Updated = updated.Format(time.RFC3339)
	return s
}

type Text struct {
	// Required elements
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

// NewText creates a new Atom text element of type 'text'.
func NewText(body string) *Text {
	return &Text{Type: "text", Body: body}
}

// NewHTML creates a new Atom text element of type 'html'.
func NewHTML(body string) *Text {
	return &Text{Type: "html", Body: body}
}

type Entry struct {
	// Required elements
	ID      string `xml:"id"`
	Title   *Text  `xml:"title"`
	Updated string `xml:"updated"`

	// Recommended elements
	Author  *Person `xml:"author,omitempty"`
	Content *Text   `xml:"content,omitempty"`
	Links   []*Link `xml:"link,omitempty"`
	Summary *Text   `xml:"summary,omitempty"`

	// Optional elements
	Categories  []*Category `xml:"category,omitempty"`
	Contributor []*Person   `xml:"contributor,omitempty"`
	Published   string      `xml:"published,omitempty"`
	Rights      *Text       `xml:"rights,omitempty"`
	Source      *Source     `xml:"source,omitempty"`
}

// NewEntry creates a new Atom entry.
func NewEntry(id string, title *Text, updated time.Time) *Entry {
	return &Entry{
		ID:      id,
		Title:   title,
		Updated: updated.Format(time.RFC3339),
	}
}

// SetAuthor sets the author of the entry.
func (e *Entry) SetAuthor(author *Person) *Entry {
	e.Author = author
	return e
}

// SetContent sets the content of the entry.
func (e *Entry) SetContent(content *Text) *Entry {
	e.Content = content
	return e
}

// AddLink adds a link to the entry.
func (e *Entry) AddLink(link *Link) *Entry {
	e.Links = append(e.Links, link)
	return e
}

// SetSummary sets a summary of the entry.
func (e *Entry) SetSummary(summary *Text) *Entry {
	e.Summary = summary
	return e
}

// AddCategory adds a category to the entry.
func (e *Entry) AddCategory(category *Category) *Entry {
	e.Categories = append(e.Categories, category)
	return e
}

// AddContributor adds a contributor to the entry.
func (e *Entry) AddContributor(contributor *Person) *Entry {
	e.Contributor = append(e.Contributor, contributor)
	return e
}

// SetPublished sets the time the entry was published.
func (e *Entry) SetPublished(published time.Time) *Entry {
	e.Published = published.Format(time.RFC3339)
	return e
}

// SetRights sets the rights of the entry.
func (e *Entry) SetRights(rights *Text) *Entry {
	e.Rights = rights
	return e
}

// SetSource sets the source of the entry.
func (e *Entry) SetSource(source *Source) *Entry {
	e.Source = source
	return e
}

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Xmlns   string   `xml:"xmlns,attr"`

	// Required elements
	ID      string `xml:"id"`
	Title   *Text  `xml:"title"`
	Updated string `xml:"updated"`

	// Strongly recommended elements
	Author *Person `xml:"author,omitempty"`
	Links  []*Link `xml:"link,omitempty"`

	// Optional elements
	Categories   []*Category `xml:"category,omitempty"`
	Contributors []*Person   `xml:"contributor,omitempty"`
	Generator    *Generator  `xml:"generator,omitempty"`
	Icon         string      `xml:"icon,omitempty"`
	Logo         string      `xml:"logo,omitempty"`
	Rights       *Text       `xml:"rights,omitempty"`
	Subtitle     *Text       `xml:"subtitle,omitempty"`

	// Entries
	Entries []*Entry `xml:"entry,omitempty"`
}

// NewFeed creates a new Atom feed.
func NewFeed(permalink string, title *Text, updated time.Time) *Feed {
	return &Feed{
		Xmlns:   "http://www.w3.org/2005/Atom",
		ID:      permalink,
		Title:   title,
		Updated: updated.Format(time.RFC3339),
	}
}

// SetAuthor sets the author of the feed.
func (f *Feed) SetAuthor(author *Person) *Feed {
	f.Author = author
	return f
}

// AddLink adds a link to the feed.
func (f *Feed) AddLink(link *Link) *Feed {
	f.Links = append(f.Links, link)
	return f
}

// AddCategory adds a category to the feed.
func (f *Feed) AddCategory(category *Category) *Feed {
	f.Categories = append(f.Categories, category)
	return f
}

// AddContributor adds a contributor to the feed.
func (f *Feed) AddContributor(contributor *Person) *Feed {
	f.Contributors = append(f.Contributors, contributor)
	return f
}

// SetGenerator sets the generator of the feed.
func (f *Feed) SetGenerator(generator *Generator) *Feed {
	f.Generator = generator
	return f
}

// SetIcon sets the icon of the feed.
// The icon identifies a small image which provides iconic visual identification for the feed.
// Icons should be square.
func (f *Feed) SetIcon(icon string) *Feed {
	f.Icon = icon
	return f
}

// SetLogo sets the logo of the feed.
// The logo identifies a larger image which provides visual identification for the feed.
// Images should be twice as wide as they are tall.
func (f *Feed) SetLogo(logo string) *Feed {
	f.Logo = logo
	return f
}

// SetRights sets the rights of the feed.
func (f *Feed) SetRights(rights *Text) *Feed {
	f.Rights = rights
	return f
}

// SetSubtitle sets the subtitle of the feed.
func (f *Feed) SetSubtitle(subtitle *Text) *Feed {
	f.Subtitle = subtitle
	return f
}

// AddEntry adds an entry to the feed.
func (f *Feed) AddEntry(entry *Entry) *Feed {
	f.Entries = append(f.Entries, entry)
	return f
}

func (f *Feed) toXML(indent bool) ([]byte, error) {
	if indent {
		// nolint: wrapcheck
		return xml.MarshalIndent(f, "", "  ")
	}
	// nolint: wrapcheck
	return xml.Marshal(f)
}

// ToXML returns the XML representation of the feed.
func (f *Feed) ToXML(indent, includeHeader bool) ([]byte, error) {
	bytes, err := f.toXML(indent)
	if err != nil {
		return nil, fmt.Errorf("error XML encoding RSS feed: %w", err)
	}
	if includeHeader {
		return append([]byte(xml.Header), bytes...), nil
	}
	return bytes, nil
}
