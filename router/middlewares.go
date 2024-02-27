package router

import (
	"context"
	"database/sql"
	"github.com/Olprog59/golog"
	"net/http"
	"time"
)

// Middleware type for convenience
type Middleware func(http.HandlerFunc) http.HandlerFunc

func dbMiddleware(db *sql.DB) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "db", db)
			next(w, r.WithContext(ctx))
		}
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		golog.Info("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		golog.Info("Completed in %v", time.Since(start))
	})
}
