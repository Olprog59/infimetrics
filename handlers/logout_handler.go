package handlers

import (
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/commons"
	"github.com/Olprog59/infimetrics/models"
	"net/http"
)

func LogoutHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session_token, err := commons.GetCookie(r, "session_token")
		if err != nil {
			golog.Warn("Could not get session token from cookie")
			return
		}

		for _, cookieName := range commons.Cookies {
			commons.ClearCookie(w, cookieName)
		}
		store, ok := models.FromContextStore(r)
		if !ok {
			golog.Warn("Could not get redis connection from context")
			return
		}
		err = store.RedisDB.Del(session_token)
		if err != nil {
			golog.Warn("Could not delete session token from redis")
			return
		}
		w.Header().Set("HX-Redirect", "/")
	}
}
