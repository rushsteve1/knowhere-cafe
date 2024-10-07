package web

import (
	"io/fs"
	"log/slog"
	"net/http"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
	"knowhere.cafe/src/shared/easy"
)

func RootHandler(staticFiles fs.FS) (out http.Handler) {
	// shh this is a secret
	mux := http.DefaultServeMux

	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/posts", http.StatusPermanentRedirect)
	})

	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		err := models.Render(r.Context(), w, "index.html", false, "")
		if err != nil {
			slog.ErrorContext(r.Context(), "index page", "error", err.Error())
		}
	})

	mux.Handle(
		"/static",
		http.StripPrefix("/static", http.FileServerFS(staticFiles)),
	)

	mux.Handle("/src", http.RedirectHandler(
		shared.REPO_URL,
		http.StatusPermanentRedirect,
	))

	mux.HandleFunc("GET /post/{id}", postHandler)

	// Apply global middleware
	out = LogMiddleware(mux)
	out = DBContextMiddleware(out)

	return out
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) == 0 {
		http.Error(w, "Missing post ID", http.StatusNotFound)
		return
	}

	ctxData := easy.Must(models.State(r.Context()))
	ctxData.DB.Find(&models.Post{}, id)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {}
