package handlers

import (
	"github.com/Olprog59/golog"
	"net/http"
)

func GetCookiesUsernameHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username, err := getCookie(r, "username")
		if err != nil {
			golog.Warn("Could not get username from cookie")
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, err = w.Write([]byte(username))
		if err != nil {
			golog.Err("Could not write username to response")
			return
		}
	}
}