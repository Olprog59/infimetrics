package router

import (
	"database/sql"
	"github.com/Olprog59/infimetrics/router/handlers"
	"net/http"
)

// Router définit l'interface pour un router HTTP, incluant les middlewares.
type Router interface {
	RegisterRoutes()
	Use(middleware func(http.Handler) http.Handler)
}

// httpRouter implémente l'interface Router.
type httpRouter struct {
	DB          *sql.DB
	middlewares []func(http.Handler) http.Handler
}

// Use ajoute un middleware à la chaîne.
func (r *httpRouter) Use(middleware func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware)
}

// NewRouter crée une nouvelle instance de Router.
func NewRouter(db *sql.DB) Router {
	return &httpRouter{
		DB: db,
	}
}

func (r *httpRouter) RegisterRoutes() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HomeHandler())

	mux.HandleFunc("/with-db", dbMiddleware(r.DB)(handlers.WithDBHandler()))

	// Appliquer les middlewares dans l'ordre inverse de leur ajout
	// pour que le premier middleware ajouté soit exécuté en premier.
	handler := http.Handler(mux)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	// Utiliser le handler final comme handler de notre serveur
	http.Handle("/", handler)
}
