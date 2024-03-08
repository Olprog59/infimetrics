package handlers

import (
	"github.com/Olprog59/infimetrics/commons"
	"net/http"
)

func AppsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		commons.RenderTemplate(w, r, "apps", &commons.Page{
			Title: "Apps",
		})
	}
}
