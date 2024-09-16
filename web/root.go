package web

import (
	"log/slog"
	"net/http"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
)

func RootHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		err := models.RenderFullTemplate(r.Context(), w, "index.html", "")
		if err != nil {
			slog.ErrorContext(r.Context(), "index page", "error", err.Error())
		}
	})

	mux.Handle("/static",
		http.StripPrefix("/static",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.FileServerFS(StaticFiles(r.Context())).ServeHTTP(w, r)
			}),
		),
	)

	mux.Handle("/src", http.RedirectHandler(
		shared.REPO_URL,
		http.StatusPermanentRedirect,
	))

	mux.HandleFunc("GET /post/{id}", postHandler)

	// Apply global middleware
	slogMux := SlogMiddleware(mux)

	return slogMux
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if len(id) == 0 {
		http.Error(w, "Missing post ID", http.StatusNotFound)
		return
	}

	ctxData := shared.Must(models.CtxState(r.Context()))
	ctxData.DB.Find(&models.Post{}, id)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {}
