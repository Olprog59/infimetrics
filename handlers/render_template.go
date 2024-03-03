package handlers

import (
	"github.com/Olprog59/golog"
	"html/template"
	"net/http"
	"path/filepath"
)

type Page struct {
	Title      string
	Data       any
	IsLoggedIn bool
	CSS        []string
	JS         []string
	Errors     []string
}

func renderTemplate(w http.ResponseWriter, r *http.Request, html string, page *Page) {
	// Chemin vers le dossier des templates
	templatesDir := "web/templates/"

	// Charge le layout principal et le template de la page demandée
	tmpl, err := template.ParseFiles(
		filepath.Join(templatesDir, "layout.html"),
		filepath.Join(templatesDir, "layouts/", html+".html"),
	)
	if err != nil {
		golog.Err("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	session_token, _ := getCookie(r, "session_token")
	//if err != nil {
	//	golog.Err("Error getting cookie: %v", err)
	//	return
	//}
	page.IsLoggedIn = session_token != ""
	// Exécute le template de layout avec les données fournies
	err = tmpl.ExecuteTemplate(w, "layout", page)
	if err != nil {
		golog.Err("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
