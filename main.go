package main

import (
	"fmt"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/commons"
	"github.com/Olprog59/infimetrics/database"
	"github.com/Olprog59/infimetrics/router"
	_ "github.com/lib/pq"
	"net/http"
)

func init() {
	golog.SetLanguage("fr")
	golog.EnableFileNameLogging()
	golog.SetTimePrecision(golog.MICRO)
	golog.SetSeparator(" | ")
}

func main() {

	golog.Success("I'm running on host %s", commons.HOST)

	db := loadDatabase()
	redis := loadRedis()
	mongo := loadMongo()

	defer func() {
		if err := db.Close(); err != nil {
			golog.Err(err.Error())
		}
	}()

	r := router.NewRouter(db, redis, mongo)
	r.Use(router.DbAndRedisMiddleware(db, redis, mongo))
	r.Use(router.AuthMiddleware)    // Ajoute le middleware d'authentification
	r.Use(router.LoggingMiddleware) // Ajoute le middleware de journalisation
	r.RegisterRoutes()

	golog.Success("Server is running on %s %s", commons.HOST, commons.PORT)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", commons.HOST, commons.PORT), nil)
	if err != nil {
		golog.Err(err.Error())
	}
}

func loadDatabase() *database.Db {
	db := database.NewDB()
	dbConnect, err := db.Connect()
	if err != nil {
		golog.Debug(err.Error())
	}
	return dbConnect
}

func loadRedis() *database.RedisDB {
	redis := database.NewRedis()
	err := redis.Ping()
	if err != nil {
		golog.Debug(err.Error())
	}
	return redis
}

func loadMongo() *database.Mongo {
	mongo := database.NewMongo()
	_, err := mongo.Connect()
	if err != nil {
		golog.Debug(err.Error())
	}
	return mongo
}
