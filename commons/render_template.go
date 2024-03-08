package commons

import (
	"github.com/Olprog59/golog"
	"html/template"
	"net/http"
	"path/filepath"
)

type Page struct {
	Title      string
	Data       any
	AppToken   string
	IsLoggedIn bool
}

var templatesDir = "web/templates/"

func RenderTemplate(w http.ResponseWriter, r *http.Request, html string, page *Page) {
	tmpl := parseTemplateFiles(w, html)
	if tmpl == nil {
		return
	}

	if page == nil {
		page = &Page{}
	}

	isLoggedin, _ := r.Context().Value("isAuthenticated").(bool)
	page.IsLoggedIn = isLoggedin

	executeTemplate(w, tmpl, "layout", page)
}

func parseTemplateFiles(w http.ResponseWriter, html string) *template.Template {
	tmpl, err := template.ParseFiles(
		filepath.Join(templatesDir, "layout.html"),
		filepath.Join(templatesDir, "layouts/", html+".html"),
	)
	if err != nil {
		golog.Err("Error parsing templates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil
	}
	return tmpl
}

func executeTemplate(w http.ResponseWriter, tmpl *template.Template, tmplName string, page *Page) {
	err := tmpl.ExecuteTemplate(w, tmplName, page)
	if err != nil {
		golog.Err("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
