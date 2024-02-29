package handlers

import (
	"net/http"
)

func PageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Set the status code to 404
	w.WriteHeader(http.StatusNotFound)

	renderTemplate(w, "404", &Page{
		Title: "Page not found",
		CSS:   []string{"404"},
	})
}
