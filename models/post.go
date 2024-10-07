package models

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Post struct {
	ModelBase
	Subject      string
	Body         string // TODO full text indexing
	RenderedBody string
	Hidden       bool

	AuthorID uuid.UUID `gorm:"not null; type:uuid"`
	Author   User

	ReplyToID *uuid.UUID `gorm:"type:uuid"`
	ReplyTo   *Post      `gorm:"foreignKey:ReplyToID"`

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
	ModelBase
	PostID uuid.UUID `gorm:"index; type:uuid"`
	UserID uuid.UUID `gorm:"index; type:uuid"`
	Reason string
}

type Vote struct {
	ModelBase
	PostID uuid.UUID `gorm:"index; type:uuid"`
	UserID uuid.UUID `gorm:"index; type:uuid"`
}
