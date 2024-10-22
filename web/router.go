package web

import (
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/rpc"

	_ "expvar"
	_ "net/http/pprof"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared/easy"
)

func checkServerError(w http.ResponseWriter, err error) bool {
	if err != nil {
		slog.Error("http response error", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

const REPO_URL string = "https://github.com/rushsteve1/knowhere-cafe"

func Router(staticFiles fs.FS) (out http.Handler) {
	// Call this just to get it setup now too
	rpc.HandleHTTP()

	// We don't use the default handler
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/search", http.StatusPermanentRedirect)
			return
		}

		// If nothing matched and we ended up here try the default mux
		// but require authentication for all of its endpoints
		RequireAuth(http.DefaultServeMux).ServeHTTP(w, r)
	})

	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServerFS(staticFiles)),
	)

	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, easy.Must(staticFiles.Open("robots.txt")))
	})

	mux.Handle(
		"/src",
		http.RedirectHandler(REPO_URL, http.StatusPermanentRedirect),
	)

	mux.Handle(
		"GET /tailscale",
		RequireAuth(http.HandlerFunc(tailscaleHandler)),
	)

	mux.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		models.NewSearch(r.URL.Query()).ServeHTTP(w, r)
	})

	mux.HandleFunc("/archive/", ArchiveHandler)

	// Apply global middleware
	return Apply(mux,
		SlogMiddleware,
		AuthMiddleware,
		DBContextMiddleware,
		GzipMiddleware,
	)
}

func tailscaleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	state, err := models.State(ctx)
	if checkServerError(w, err) {
		return
	}

	lc, err := state.Tsnet.LocalClient()
	if checkServerError(w, err) {
		return
	}

	status, err := lc.Status(ctx)
	if checkServerError(w, err) {
		return
	}

	status.WriteHTML(w)
}
