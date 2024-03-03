package router

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/database"
	"net/http"
	"strings"
	"time"
)

const (
	SignInPath  = "/sign-in"
	SignUpPath  = "/sign-up"
	FaviconPath = "/favicon.ico"
	StaticPath  = "/static/"
)

func isPublicPath(path string) bool {
	publicPaths := []string{SignInPath, SignUpPath, FaviconPath, StaticPath}
	for _, p := range publicPaths {
		if path == p || strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, StaticPath) {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		golog.Debug("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		golog.Debug("Completed in %v", time.Since(start))
	})
}

func DbAndRedisMiddleware(db *sql.DB, redis *database.RedisDB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), database.DbKey, db)
			ctx = context.WithValue(ctx, database.RedisKey, redis)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		username, ok := isAuthenticated(r)
		if !ok {
			golog.Warn("User is not authenticated")
			golog.Debug("Request: %s", r.URL.Path)

			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusUnauthorized)
				_, err := fmt.Fprint(w, "Unauthorized - Please log in")
				if err != nil {
					return
				}
				return
			}
			w.Header().Set("HX-Redirect", SignInPath)
			// Utilise HX-Redirect pour les redirections côté client avec HTMX
			http.Redirect(w, r, SignInPath, http.StatusSeeOther)
			return
		}
		w.Header().Set("HX-Current-Username", username)
		next.ServeHTTP(w, r)
	})
}

func isAuthenticated(r *http.Request) (string, bool) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		golog.Warn("Error getting cookie")
		return "", false
	}
	redis, ok := database.FromContextRedis(r)
	if !ok {
		golog.Warn("Could not get Redis connection from context")
		return "", false
	}

	_, err = redis.HGet(sessionToken.Value, "userID")
	if err != nil {
		golog.Warn("Error getting value from Redis")
		return "", false
	}

	username, err := redis.HGet(sessionToken.Value, "username")
	if err != nil {
		golog.Warn("Error getting value from Redis")
		return "", false
	}

	//golog.Success("User %s is authenticated", username)
	return username, true
}
