package handlers

import (
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/database"
	"net/http"
)

func LogoutHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session_token, err := getCookie(r, "session_token")
		if err != nil {
			golog.Warn("Could not get session token from cookie")
			return
		}
		clearCookieHandler(w, "session_token")
		redis, ok := database.FromContextRedis(r)
		if !ok {
			golog.Warn("Could not get redis connection from context")
			return
		}
		err = redis.Del(session_token)
		if err != nil {
			golog.Warn("Could not delete session token from redis")
			return
		}
		w.Header().Set("HX-Redirect", "/")
	}
}
