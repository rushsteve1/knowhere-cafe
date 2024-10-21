package models

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"gorm.io/gorm"
)

type SitePrefs struct {
	gorm.Model
	Domain  string
	Score   int
	Blocked bool
}

type Search struct {
	gorm.Model
	Terms   string
	Results []SearchResult `gorm:"type:jsonb;serializer:json"`
}

func NewSearch(params url.Values) Search {
	return Search{
		Terms: params.Get("terms"),
	}
}

func (s Search) TemplateName() string { return "search.html" }

func (s Search) TitleString() string        { return fmt.Sprintf("Search for \"%s\"", s.Terms) }
func (s Search) BodyString() string         { return "" }
func (s Search) PublishedAt() time.Time     { return time.Now() }
func (s Search) Markdown(w io.Writer) error { return simpleMarkdown(w, s) }
func (s Search) Etag() string               { return timeHash(&s.UpdatedAt) }
func (s Search) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	formatRenderHandler(w, r, s.TemplateName(), s)
}

type SearchResult struct {
	gorm.Model
	URL     string
	Title   string
	Summary string
}
