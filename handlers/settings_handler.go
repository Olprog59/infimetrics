package handlers

import (
	"github.com/Olprog59/infimetrics/commons"
	"net/http"
)

func SettingsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		commons.RenderTemplate(w, r, "settings", &commons.Page{
			Title: "Settings",
		})
	}
}
