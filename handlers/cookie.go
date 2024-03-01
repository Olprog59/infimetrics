package handlers

import (
	"github.com/Olprog59/infimetrics/commons"
	"net/http"
	"time"
)

func setCookieHandler(w http.ResponseWriter, r *http.Request, value string) {
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(commons.TimeoutCookie),
		MaxAge:   int(commons.TimeoutCookie.Seconds()),
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}

func clearCookieHandler(w http.ResponseWriter, name string) {
	cookie := http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-commons.TimeoutCookie),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}

func getCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
