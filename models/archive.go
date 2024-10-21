package models

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-shiori/go-readability"
	"github.com/go-shiori/obelisk"
	"gorm.io/gorm"
)

type Archive struct {
	gorm.Model
	Title         string
	Byline        string
	Content       string
	TextContent   string
	Length        int
	Excerpt       string
	SiteName      string
	Image         string
	Favicon       string
	PublishedTime *time.Time
	Read          bool
	URL           string `gorm:"unique"`
	HTML          sql.NullString
	Notes         []Note `gorm:"many2many:archive_notes"`
}

func NewArchive(ctx context.Context, u *url.URL) (out Archive, err error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return out, err
	}

	arc := obelisk.Archiver{}
	doc, _, err := arc.Archive(ctx, obelisk.Request{Input: resp.Body})
	if err != nil {
		return out, err
	}

	article, err := readability.FromReader(resp.Body, u)
	if err != nil {
		return out, err
	}

	return Archive{
		URL:           u.String(),
		Title:         article.Title,
		Byline:        article.Byline,
		Content:       article.Content,
		TextContent:   article.TextContent,
		Length:        article.Length,
		Excerpt:       article.Excerpt,
		SiteName:      article.SiteName,
		Image:         article.Image,
		Favicon:       article.Favicon,
		PublishedTime: article.PublishedTime,
		HTML:          sql.NullString{String: string(doc), Valid: true},
	}, nil
}

func (a Archive) TemplateName() string       { return "archive.html" }
func (a Archive) TitleString() string        { return a.Title }
func (a Archive) BodyString() string         { return a.Content }
func (a Archive) PublishedAt() time.Time     { return a.UpdatedAt }
func (a Archive) Markdown(w io.Writer) error { return simpleMarkdown(w, a) }
func (a Archive) Etag() string               { return timeHash(&a.UpdatedAt) }
func (a Archive) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	formatRenderHandler(w, r, a.TemplateName(), a)
}

type ArchiveList struct {
	Page    int
	List    []Archive
	Current *Archive
}

func (al ArchiveList) TemplateName() string { return "archive.html" }
func (al ArchiveList) TitleString() string { return fmt.Sprintf("Archives page %d", al.Page) }
func (al ArchiveList) BodyString() string  { return "" }
func (al ArchiveList) PublishedAt() time.Time {
	if al.Current != nil {
		return al.Current.PublishedAt()
	}
	if len(al.List) > 0 {
		return al.List[0].PublishedAt()
	}
	return time.Now()
}
func (al ArchiveList) Markdown(
	w io.Writer,
) error {
	return simpleMarkdown(w, al)
}
func (al ArchiveList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	formatRenderHandler(w, r, al.TemplateName(), al)
}
func (al ArchiveList) Etag() string {
	t := al.PublishedAt()
	return timeHash(&t)
}
