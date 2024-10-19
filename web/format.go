package web

import (
	"cmp"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"log/slog"
	"mime"
	"net/http"
	"strings"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared/easy"
)

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

func formatRenderer(
	w http.ResponseWriter,
	r *http.Request,
	name string,
	data models.Renderable,
) {
	ctx := r.Context()
	state := easy.Must(models.State(ctx))

	format := r.URL.Query().Get(FORMAT_KEY)
	accept := r.Header.Get(ACCEPT_HEADER)
	format = cmp.Or(format, accept)

	target := r.Header.Get(UP_TARGET_HEADER)

	slog.Debug("response format", "format", format)

	if containsOne(format, "html", html_mime) {
		w.Header().Set(CONTENT_TYPE_HEADER, html_mime)
		easy.Expect(state.Templ.Render(w, name, target, state.Flags.Dev, data))
	} else if containsOne(format, "js", "json", json_mime) {
		w.Header().Set(CONTENT_TYPE_HEADER, json_mime)
		easy.Expect(json.NewEncoder(w).Encode(data))
	} else if containsOne(format, "md", "markdown", "txt", "text", md_mime, text_mime) {
		w.Header().Set(CONTENT_TYPE_HEADER, md_mime)
		easy.Expect(data.Markdown(w))
	} else if containsOne(format, "xml", xml_mime) {
		w.Header().Set(CONTENT_TYPE_HEADER, xml_mime)
		easy.Expect(xml.NewEncoder(w).Encode(data))
	} else if containsOne(format, "go", "gob", gob_mime) {
		w.Header().Set(CONTENT_TYPE_HEADER, gob_mime)
		easy.Expect(gob.NewEncoder(w).Encode(data))
	} else {
		http.Error(w, "unknown format", http.StatusNotAcceptable)
	}
}

func containsOne(s string, substrs ...string) bool {
	for _, n := range substrs {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}
