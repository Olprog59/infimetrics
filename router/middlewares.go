package router

import (
	"context"
	"database/sql"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/commons"
	"github.com/Olprog59/infimetrics/database"
	"net/http"
	"strings"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/") {
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
		// Check if user is authenticated
		if r.URL.Path == "/sign-in" ||
			strings.HasPrefix(r.URL.Path, "/sign-up") ||
			r.URL.Path == "/favicon.ico" ||
			strings.HasPrefix(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}
		if !isAuthenticated(r) {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isAuthenticated(r *http.Request) bool {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		golog.Warn("Error getting cookie")
		return false
	}
	redis, ok := commons.FromContextRedis(r)
	if !ok {
		golog.Warn("Could not get Redis connection from context")
		return false
	}

	username, err := redis.Get(sessionToken.Value)
	if err != nil {
		golog.Warn("Error getting value from Redis")
		return false
	}

	golog.Success("User %s is authenticated", username)
	return true
}
