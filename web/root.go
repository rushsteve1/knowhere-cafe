package web

import (
	"io/fs"
	"net/http"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
	"knowhere.cafe/src/shared/log"
)

func RootHandler(staticFiles fs.FS) http.Handler {
	http.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		err := models.Render(r.Context(), w, "index.html", "")
		if err != nil {
			log.ErrorContext(r.Context(), "index page", "error", err.Error())
		}
	})

	http.Handle("/static", http.StripPrefix("/static", http.FileServerFS(staticFiles)))

	http.Handle("/src", http.RedirectHandler(
		shared.REPO_URL,
		http.StatusPermanentRedirect,
	))

	http.HandleFunc("GET /post/{id}", postHandler)

	// Apply global middleware
	logMux := LogMiddleware(http.DefaultServeMux)

	return logMux
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) == 0 {
		http.Error(w, "Missing post ID", http.StatusNotFound)
		return
	}

	ctxData := log.Must(models.CtxState(r.Context()))
	ctxData.DB.Find(&models.Post{}, id)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {}
