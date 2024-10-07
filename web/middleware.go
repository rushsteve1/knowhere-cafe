package web

import (
	"context"
	"log/slog"
	"net/http"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
	"knowhere.cafe/src/shared/easy"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.InfoContext(
			r.Context(),
			"http request",
			"method", r.Method,
			"url", r.URL.String(),
			"peer", r.RemoteAddr,
		)

		next.ServeHTTP(w, r)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return nil
}

func PermissionsMiddleware(
	next http.Handler,
	perm models.Permissions,
) http.Handler {
	return nil
}

func DBContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		state := easy.Must(models.State(ctx))
		state.DB = state.DB.WithContext(ctx)
		ctx = context.WithValue(ctx, shared.CTX_STATE_KEY, state)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
