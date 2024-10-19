package web

import (
	"io"
	"io/fs"
	"net/http"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
	"knowhere.cafe/src/shared/easy"
)

func Router(staticFiles fs.FS) (out http.Handler) {
	// shh this is a secret
	mux := http.DefaultServeMux

	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/search", http.StatusPermanentRedirect)
	})

	mux.Handle(
		"/static",
		http.StripPrefix("/static", http.FileServerFS(staticFiles)),
	)

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, easy.Must(staticFiles.Open("robots.txt")))
	})

	mux.Handle("/src", http.RedirectHandler(
		shared.REPO_URL,
		http.StatusPermanentRedirect,
	))

	mux.HandleFunc("/search", searchHandler)

	// Apply global middleware
	out = SlogMiddleware(mux)
	out = DBContextMiddleware(out)

	return out
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	formatRenderer(w, r, "search.html", models.Search{})
}
