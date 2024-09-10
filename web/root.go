package web

import (
	"net/http"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
)

func RootHandler() http.Handler {
	mux := http.NewServeMux()

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

func usersHandler(w http.)