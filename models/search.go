package models

import (
	"fmt"
	"io"
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
	Terms string
}

func (s Search) Title() string { return fmt.Sprintf("Search for \"%s\"", s.Terms) }
func (s Search) Body() string { return "" }
func (s Search) PublishedAt() time.Time { return time.Now() }
func (s Search) Markdown(w io.Writer) error {
	_, err := io.WriteString(w, "# " + s.Title())
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, s.PublishedAt().Format(time.RFC3339))
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, s.Body())
	if err != nil {
		return err
	}
	return nil
}