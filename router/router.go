package router

import (
	"database/sql"
	"github.com/Olprog59/infimetrics/database"
	"github.com/Olprog59/infimetrics/handlers"
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
	Redis       *database.RedisDB
	middlewares []func(http.Handler) http.Handler
}

// Use ajoute un middleware à la chaîne.
func (r *httpRouter) Use(middleware func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware)
}

// NewRouter crée une nouvelle instance de Router.
func NewRouter(db *sql.DB, redis *database.RedisDB) Router {
	return &httpRouter{
		DB:    db,
		Redis: redis,
	}
}

func (r *httpRouter) RegisterRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.PageNotFoundHandler(w, r)
	})

	// Routes pour les fichiers statiques
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/assets"))))

	// Routes pour les pages
	mux.HandleFunc("GET /{$}", handlers.HomeHandler())

	// Routes pour la connexion et l'inscription
	mux.HandleFunc("/sign-in", handlers.LoginHandler())
	mux.HandleFunc("/sign-up", handlers.RegisterHandler())

	// Routes pour la vérification des données d'inscription
	mux.HandleFunc("POST /sign-up/email", handlers.SignUpEmail())
	mux.HandleFunc("POST /sign-up/username", handlers.SignUpUsername())
	mux.HandleFunc("POST /sign-up/password", handlers.SignUpPassword())
	mux.HandleFunc("POST /sign-up/same-password", handlers.SignUpPasswordSame())

	mux.HandleFunc("/dashboard", handlers.DashboardHandler())

	handler := http.Handler(mux)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	http.Handle("/", handler)
}
