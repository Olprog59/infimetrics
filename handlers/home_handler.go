package handlers

import (
	"github.com/Olprog59/infimetrics/commons"
	"net/http"
)

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		session_token, err := commons.GetCookie(r, "session_token")
		if err != nil {
			return
		}
		commons.RenderTemplate(w, r, "dashboard", &commons.Page{
			Title:      "Dashboard",
			IsLoggedIn: session_token != "",
		})
	}
}
