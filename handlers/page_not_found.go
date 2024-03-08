package handlers

import (
	"github.com/Olprog59/infimetrics/commons"
	"net/http"
)

func PageNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// Set the status code to 404
	w.WriteHeader(http.StatusNotFound)

	commons.RenderTemplate(w, r, "404", &commons.Page{
		Title: "Page not found",
	})
}
