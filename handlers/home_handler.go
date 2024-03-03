package handlers

import (
	"net/http"
)

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// template
		renderTemplate(w, "dashboard", &Page{
			Title: "Dashboard",
		})
	}
}
