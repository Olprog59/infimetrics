package handlers

import (
	"github.com/Olprog59/infimetrics/commons"
	"net/http"
	"time"
)

type CookieName string

const (
	SessionToken CookieName = "session_token"
	Username     CookieName = "username"
)

func ListCookieNames() []CookieName {
	return []CookieName{SessionToken, Username}
}

func setCookieHandler(w http.ResponseWriter, r *http.Request, key CookieName, value string) {
	cookie := http.Cookie{
		Name:     string(key),
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(commons.TimeoutCookie),
		MaxAge:   int(commons.TimeoutCookie.Seconds()),
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}

func clearCookieHandler(w http.ResponseWriter, name CookieName) {
	cookie := http.Cookie{
		Name:     string(name),
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
