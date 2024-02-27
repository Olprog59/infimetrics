package main

import (
	"context"
	"fmt"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/appconfig"
	"github.com/Olprog59/infimetrics/database"
	router "github.com/Olprog59/infimetrics/router"
	_ "github.com/lib/pq"
	"net/http"
)

func init() {
	golog.SetLanguage("fr")
	golog.EnableFileNameLogging()
}

func main() {
	cfg, err := appconfig.LoadFromPath(context.Background(), "pkl/int/appConfig.pkl")
	if err != nil {
		golog.Err(err.Error())
		panic(err)
	}
	golog.Info("I'm running on host %s", cfg.Host)

	db := database.NewDB(cfg.Database)
	dbConnect, err := db.Connect()
	if err != nil {
		golog.Debug(err.Error())
	}

	defer func() {
		if err := db.Close(); err != nil {
			golog.Err(err.Error())
		}
	}()

	r := router.NewRouter(dbConnect)
	r.Use(router.LoggingMiddleware) // Ajoute le middleware de journalisation
	r.RegisterRoutes()

	golog.Info("Server is running on %s %d", cfg.Host, cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil)
	if err != nil {
		golog.Err(err.Error())
	}
}
