package web

import (
	"compress/gzip"
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"knowhere.cafe/src/models"
	"knowhere.cafe/src/shared"
	"knowhere.cafe/src/shared/easy"
)

type Middleware = func(next http.Handler) http.Handler

func Apply(top http.Handler, chain ...Middleware) http.Handler {
	cur := top
	for _, fn := range chain {
		cur = fn(cur)
	}
	return cur
}

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

func DBContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		state := easy.Must(models.State(ctx))
		state.DB = state.DB.WithContext(ctx)
		ctx = context.WithValue(ctx, shared.CTX_STATE_KEY, state)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// https://gist.github.com/the42/1956518
type gzipWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

const ACCEPT_ENCODING = "Accept-Encoding"
const CONTENT_ENCODING = "Content-Encoding"

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get(ACCEPT_ENCODING), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set(CONTENT_ENCODING, "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		next.ServeHTTP(gzipWriter{gz, w}, r)
	})
}
