package models

import (
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Subject      string
	Body         string // TODO full text indexing
	RenderedBody string
	Hidden       bool

	AuthorID uint
	Author   User

	ReplyToID *uint
	ReplyTo   *Post `gorm:"foreignKey:ReplyToID"`

	Children []Post `gorm:"foreignKey:ReplyToID"`
	Reports  []Report
	Votes    []Vote
}

// Posts implement the ServeHTTP interface allowing them to mux.
func (p Post) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	// Strip off the post's path prefix if it wasn't already
	prefix := fmt.Sprintf("/posts/%d", p.ID)
	stripped := http.StripPrefix(prefix, mux)
	stripped.ServeHTTP(w, r)
}

type Report struct {
	gorm.Model
	PostID uint `gorm:"index"`
	UserID uint `gorm:"index"`
	Reason string
}

type Vote struct {
	gorm.Model
	PostID uint `gorm:"index"`
	UserID uint `gorm:"index"`
}
