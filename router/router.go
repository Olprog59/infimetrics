package router

import (
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/database"
	"github.com/Olprog59/infimetrics/handlers"
	"github.com/Olprog59/infimetrics/handlers/api"
	"net/http"
)

// Router définit l'interface pour un router HTTP, incluant les middlewares.
type Router interface {
	RegisterRoutes()
	Use(middleware func(http.Handler) http.Handler)
}

// httpRouter implémente l'interface Router.
type httpRouter struct {
	DB          *database.Db
	Redis       *database.RedisDB
	Mongo       *database.Mongo
	middlewares []func(http.Handler) http.Handler
}

// Use ajoute un middleware à la chaîne.
func (r *httpRouter) Use(middleware func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware)
}

// NewRouter crée une nouvelle instance de Router.
func NewRouter(db *database.Db, redis *database.RedisDB, mongo *database.Mongo) Router {
	return &httpRouter{
		DB:    db,
		Redis: redis,
		Mongo: mongo,
	}
}

func (r *httpRouter) RegisterRoutes() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		golog.Debug("Page not found: %s", r.URL.Path)
		handlers.PageNotFoundHandler(w, r)
	})

	// Routes pour les fichiers statiques
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/assets"))))

	// Routes pour les pages
	mux.HandleFunc("GET /{$}", handlers.HomeHandler())

	// Routes pour la connexion et l'inscription
	mux.HandleFunc("/sign-in", handlers.LoginHandler())
	mux.HandleFunc("/sign-up", handlers.RegisterHandler())
	mux.HandleFunc("/logout", handlers.LogoutHandler())

	// Routes pour la vérification des données d'inscription
	mux.HandleFunc("POST /sign-up/email", handlers.SignUpEmail())
	mux.HandleFunc("POST /sign-up/username", handlers.SignUpUsername())
	mux.HandleFunc("POST /sign-up/password", handlers.SignUpPassword())
	mux.HandleFunc("POST /sign-up/same-password", handlers.SignUpPasswordSame())

	mux.HandleFunc("GET /apps", handlers.AppsHandler())

	// Routes pour les applications
	mux.HandleFunc("GET /api/v1/apps", api.ApiAppsHandler())
	mux.HandleFunc("GET /api/v1/app", api.ApiNewAppsHandler())
	mux.HandleFunc("POST /api/v1/app", api.ApiNewAppsPostHandler())
	mux.HandleFunc("DELETE /api/v1/app/{token}", api.ApiDeleteAppsPostHandler())
	mux.HandleFunc("DELETE /api/v1/app/modal", api.ApiNewAppsDeleteHandler())

	// Routes pour les logs
	mux.HandleFunc("GET /api/v1/apps/{token}", api.ApiOneAppHandler())
	mux.HandleFunc("GET /api/v1/logs/{token}", api.ApiLogsHandler())
	//mux.HandleFunc("GET /api/v1/logs/{token}/refresh", api.ApiLogsWatchHandler())

	mux.HandleFunc("/settings", handlers.SettingsHandler())

	handler := http.Handler(mux)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	http.Handle("/", handler)
}
