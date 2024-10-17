package web

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
	"knowhere.cafe/src/shared/easy"
)

func userHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	state := easy.Must(models.State(ctx))

	id := r.PathValue("id")
	easy.Assert(len(id) > 0)

	var user models.User
	res := state.DB.Find(&user, id)
	easy.Check(res.Error)

	// Check the ETag to see if we have to render any of this or can just return
	etag := user.ETag()
	if etag == r.Header.Get(shared.NON_MATCH_HEADER) {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set(shared.ETAG_HEADER, etag)

	// Check the Unpoly target an set the Vary header
	target := r.Header.Get(shared.UP_TARGET_HEADER)
	if len(target) > 0 {
		w.Header().Add(shared.VARY_HEADER, target)
	}

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/users/{id}/{$}",
		func(w http.ResponseWriter, r *http.Request) {
			models.Render(ctx, w, "user.html", target == "main", user)
		},
	)

	mux.HandleFunc(
		"/users/{id}/json",
		func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(user)
		},
	)

	mux.HandleFunc(
		"/users/{id}/xml",
		func(w http.ResponseWriter, r *http.Request) {
			xml.NewEncoder(w).Encode(user)
		},
	)

	mux.HandleFunc(
		"/users/{id}/popup",
		func(w http.ResponseWriter, r *http.Request) {
			if target != "popup" {
				http.Redirect(w, r, "/users/"+user.ID.String(), http.StatusTemporaryRedirect)
			}
			models.Render(ctx, w, "user_popup.html", true, user)
		},
	)

	mux.ServeHTTP(w, r)
}