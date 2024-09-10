package models

import (
	"database/sql"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ReplyTo      sql.NullInt64 `gorm:"references:posts(id)"`
	Subject      string
	Body         string // TODO full text indexing
	RenderedBody string
	Points       int
}

type Report struct {
	PostID uint `gorm:"references:posts(id)"`
	UserID uint `gorm:"references:users(id)"`
	Reason string
}

// Posts implement the ServeHTTP interface allowing them to mux.
func (p Post) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()

	// Strip off the post's path prefix if it wasn't already
	prefix := fmt.Sprintf("/posts/%d", p.ID)
	stripped := http.StripPrefix(prefix, mux)
	stripped.ServeHTTP(w, r)
}

func (p Post) QueryChildren() []Post {
	return nil
}
