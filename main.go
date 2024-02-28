package main

import (
	"context"
	"database/sql"
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
	cfg, err := loadConfig()
	if err != nil {
		golog.Err(err.Error())
		panic(err)
	}
	golog.Info("I'm running on host %s", cfg.Host)

	db, dbConnect := loadDatabase(cfg, err)
	redis := loadRedis(cfg)

	defer func() {
		if err := db.Close(); err != nil {
			golog.Err(err.Error())
		}
	}()

	r := router.NewRouter(dbConnect, redis)
	r.Use(router.LoggingMiddleware) // Ajoute le middleware de journalisation
	r.RegisterRoutes()

	golog.Info("Server is running on %s %d", cfg.Host, cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil)
	if err != nil {
		golog.Err(err.Error())
	}
}

func loadDatabase(cfg *appconfig.AppConfig, err error) (database.IDB, *sql.DB) {
	db := database.NewDB(cfg.Database)
	dbConnect, err := db.Connect()
	if err != nil {
		golog.Debug(err.Error())
	}
	return db, dbConnect
}

func loadRedis(cfg *appconfig.AppConfig) *database.RedisDB {
	redis := database.NewRedis(cfg.Redis)
	err := redis.Ping()
	if err != nil {
		golog.Debug(err.Error())
	}
	return redis
}

func loadConfig() (*appconfig.AppConfig, error) {
	return appconfig.LoadFromPath(context.Background(), "pkl/int/appConfig.pkl")
}
