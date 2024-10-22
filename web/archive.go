package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"gorm.io/gorm/clause"
	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
	"knowhere.cafe/src/shared/easy"
)

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
				http.StatusCreated,
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

	mux.Handle(
		"PATCH /archive/{id}",
		RequireAuth(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				state := easy.Must(models.State(r.Context()))

				idstr := r.PathValue("id")
				id, err := strconv.Atoi(idstr)
				if checkServerError(w, err) {
					return
				}

				err = r.ParseForm()
				if checkServerError(w, err) {
					return
				}

				read, err := strconv.ParseBool(r.Form.Get("read"))
				if checkServerError(w, err) {
					return
				}

				var ar models.Archive
				ar.ID = uint(id)
				ar.Read = read
				res := state.DB.Save(&ar)
				if checkServerError(w, res.Error) {
					return
				}

				w.WriteHeader(http.StatusNoContent)
			}),
		),
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
