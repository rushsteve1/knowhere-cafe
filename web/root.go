package web

import (
	"net/http"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
)

func RootHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		// TODO
		w.Write([]byte("hello world"))
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/static",
		http.StripPrefix("/static", http.FileServerFS(StaticFiles)),
	)

	mux.Handle("/src", http.RedirectHandler(
		shared.REPO_URL,
		http.StatusPermanentRedirect,
	))

	mux.HandleFunc("GET /post/{id}", postsHandler)

	// Apply global middleware
	slogMux := SlogMiddleware(mux)

	return slogMux
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) == 0 {
		http.Error(w, "Missing post ID", http.StatusNotFound)
		return
	}

	ctxData := shared.CtxData(r.Context())
	ctxData.DB.Find(&models.Post{}, id)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {}
