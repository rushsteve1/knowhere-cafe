package web

import (
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"gorm.io/gorm/clause"
	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
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
	// shh this is a secret
	mux := http.DefaultServeMux

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/search", http.StatusPermanentRedirect)
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

	mux.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		models.NewSearch(r.URL.Query()).ServeHTTP(w, r)
	})

	mux.HandleFunc("/archive/", ArchiveHandler)

	// Apply global middleware
	return Apply(mux,
		SlogMiddleware,
		GzipMiddleware,
		DBContextMiddleware,
	)
}

func ArchiveHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	state := easy.Must(models.State(ctx))

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 0
	}

	al := models.ArchiveList{Page: page}
	res := state.DB.Limit(shared.LIMIT).
		Offset(page * shared.LIMIT).
		Order("updated_at DESC").
		Find(&al.List)
	if checkServerError(w, res.Error) {
		return
	}

	mux := http.NewServeMux()

	mux.Handle("GET /archive/{$}", al)

	mux.HandleFunc(
		"POST /archive/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			formUrl, err := url.Parse(r.FormValue("url"))
			if checkServerError(w, err) {
				return
			}

			archive, err := models.NewArchive(r.Context(), formUrl)
			if checkServerError(w, err) {
				return
			}

			slog.Debug("inserting new archive", "URL", archive.URL)
			res = state.DB.Clauses(clause.OnConflict{UpdateAll: true}).
				Create(&archive)
			if checkServerError(w, res.Error) {
				return
			}

			http.Redirect(
				w,
				r,
				fmt.Sprintf("/archive/%d", archive.ID),
				http.StatusSeeOther,
			)
		},
	)

	mux.HandleFunc(
		"GET /archive/{id}",
		func(w http.ResponseWriter, r *http.Request) {
			idstr := r.PathValue("id")
			id, err := strconv.Atoi(idstr)
			if checkServerError(w, err) {
				return
			}

			res := state.DB.First(&al.Current, id)
			if checkServerError(w, res.Error) {
				return
			}

			al.ServeHTTP(w, r)
		},
	)

	mux.HandleFunc(
		"GET /archive/{id}/html",
		func(w http.ResponseWriter, r *http.Request) {
			idstr := r.PathValue("id")
			id, err := strconv.Atoi(idstr)
			if checkServerError(w, err) {
				return
			}

			var a models.Archive
			res := state.DB.First(&a, id)
			if checkServerError(w, res.Error) {
				return
			}

			w.Write([]byte(a.HTML))
		},
	)

	mux.ServeHTTP(w, r)
}
