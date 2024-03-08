package commons

import (
	"net/http"
	"time"
)

type CookieName string

const (
	SessionToken CookieName = "session_token"
	Username     CookieName = "username"
)

var Cookies []CookieName = []CookieName{SessionToken, Username}

func SetCookie(w http.ResponseWriter, r *http.Request, key CookieName, value string) {
	cookie := http.Cookie{
		Name:     string(key),
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(TimeoutCookie),
		MaxAge:   int(TimeoutCookie.Seconds()),
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}

func ClearCookie(w http.ResponseWriter, name CookieName) {
	cookie := http.Cookie{
		Name:     string(name),
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-TimeoutCookie),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}

func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
