package handlers

import (
	"net/http"
)

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		session_token, err := getCookie(r, "session_token")
		if err != nil {
			return
		}
		renderTemplate(w, "dashboard", &Page{
			Title:      "Dashboard",
			CSS:        []string{"dashboard"},
			IsLoggedIn: session_token != "",
		})
	}
}
