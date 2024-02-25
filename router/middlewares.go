package router

import (
	"context"
	"database/sql"
	"log"
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
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed in %v", time.Since(start))
	})
}
