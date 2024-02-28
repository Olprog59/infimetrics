package handlers

import (
	"html/template"
	"net/http"
)

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// template
		tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			return
		}
	}
}
