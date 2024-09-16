package web

import (
	"log/slog"
	"net/http"

	"knowhere.cafe/src/models"
)

func SlogMiddleware(next http.Handler) http.Handler {
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
