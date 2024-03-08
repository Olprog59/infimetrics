package router

import (
	"context"
	"fmt"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/database"
	"github.com/Olprog59/infimetrics/models"
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

var publicPaths = []string{SignInPath, SignUpPath, FaviconPath, StaticPath}

func isPublicPath(path string) bool {
	for _, p := range publicPaths {
		if path == p || strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isPublicPath(r.URL.Path) {
			start := time.Now()
			defer golog.Debug("Completed in %v", time.Since(start))
			golog.Debug("Started %s %s", r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

func DbAndRedisMiddleware(db *database.Db, redis *database.RedisDB, mongo *database.Mongo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			store := &models.Store{
				Db:      db,
				RedisDB: redis,
				Mongo:   mongo,
			}
			ctx := context.WithValue(r.Context(), models.StoreKey, store)
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

		newCtx, username, ok := isAuthenticated(r)

		if !ok {
			handleUnauthenticatedUser(w, r)
			return
		}

		w.Header().Set("HX-Current-Username", username)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

func handleUnauthenticatedUser(w http.ResponseWriter, r *http.Request) {
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
	clearSessionCookie(w)
	http.Redirect(w, r, SignInPath, http.StatusSeeOther)
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func isAuthenticated(r *http.Request) (context.Context, string, bool) {
	newCtx := context.WithValue(r.Context(), "isAuthenticated", false)
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		golog.Warn("Error getting cookie")
		return newCtx, "", false
	}

	store, ok := models.FromContextStore(r)
	if !ok {
		golog.Warn("Could not get Store connection from context")
		return newCtx, "", false
	}

	_, err = store.RedisDB.HGet(sessionToken.Value, "userID")
	if err != nil {
		golog.Warn("Error getting value from Redis")
		return newCtx, "", false
	}

	username, err := store.RedisDB.HGet(sessionToken.Value, "username")
	if err != nil {
		golog.Warn("Error getting value from Redis")
		return newCtx, "", false
	}

	newCtx = context.WithValue(r.Context(), "isAuthenticated", true)
	return newCtx, username, true
}
