package models

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"strings"
	"time"

	"knowhere.cafe/src/shared/easy"
)

// Renderable describes an entity that can be rendered over HTTP using the
// [[formatRenderer]] function to HTML and Markdown.
type Renderable interface {
	http.Handler
	TemplateName() string
	TitleString() string
	BodyString() string
	PublishedAt() time.Time
	Markdown(w io.Writer) error
	Etag() string
}

const FORMAT_KEY = "format"
const UP_TARGET_HEADER = "X-Up-Target"
const ACCEPT_HEADER = "Accept"
const CONTENT_TYPE_HEADER = "Content-Type"

// Pre-defining these because it's easier
var (
	html_mime = mime.TypeByExtension(".html")
	json_mime = mime.TypeByExtension(".json")
	xml_mime  = mime.TypeByExtension(".xml")
	md_mime   = mime.TypeByExtension(".md")
	text_mime = mime.TypeByExtension(".txt")
	gob_mime  = mime.TypeByExtension(".gob")
)

// Format and respond the page according to the request's `Accept` header
// This is usually the end-of-the-line function for a GET request
// after data is pulled out of the database.
func formatRenderHandler(
	w http.ResponseWriter,
	r *http.Request,
	name string,
	data Renderable,
) {
	ctx := r.Context()
	state := easy.Must(State(ctx))

	format := r.Header.Get(ACCEPT_HEADER)
	slog.Debug("response format", "format", format)
	target := r.Header.Get(UP_TARGET_HEADER)
	slog.Debug("response target", "target", target)

	formats := strings.Split(format, ",")

	// Go through every possible value in request order
	// ignoring the quality because
	for _, f := range formats {
		if prefixAny(f, html_mime) {
			w.Header().Set(CONTENT_TYPE_HEADER, html_mime)
			easy.Expect(
				state.Templ.Render(w, name, target, state.Flags.Dev, data),
			)
			return
		} else if prefixAny(f, json_mime) {
			w.Header().Set(CONTENT_TYPE_HEADER, json_mime)
			easy.Expect(json.NewEncoder(w).Encode(data))
			return
		} else if prefixAny(f, md_mime, text_mime) {
			w.Header().Set(CONTENT_TYPE_HEADER, md_mime)
			easy.Expect(data.Markdown(w))
			return
		} else if prefixAny(f, xml_mime) {
			w.Header().Set(CONTENT_TYPE_HEADER, xml_mime)
			easy.Expect(xml.NewEncoder(w).Encode(data))
			return
		} else if prefixAny(f, gob_mime) {
			w.Header().Set(CONTENT_TYPE_HEADER, gob_mime)
			easy.Expect(gob.NewEncoder(w).Encode(data))
			return
		}
	}
	http.Error(w, "unknown format", http.StatusNotAcceptable)
}

func prefixAny(prefix string, substrs ...string) bool {
	for _, s := range substrs {
		if strings.HasPrefix(prefix, s) {
			return true
		}
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func timeHash(t *time.Time) string {
	sum := md5.Sum([]byte(t.Format(time.RFC3339)))
	return base64.StdEncoding.EncodeToString(sum[:])
}

func simpleMarkdown(w io.Writer, r Renderable) error {
	_, err := io.WriteString(w, "# "+r.TitleString()+"\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, r.PublishedAt().Format(time.RFC3339)+"\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, r.BodyString())
	if err != nil {
		return err
	}
	return nil
}
