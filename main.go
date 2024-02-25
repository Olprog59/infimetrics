package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/olprog59/infimetrics/commons"
	"github.com/olprog59/infimetrics/database"
	router "github.com/olprog59/infimetrics/router"
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
}

func main() {
	db := database.NewDB()
	dbConnect, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	r := router.NewRouter(dbConnect)
	r.Use(router.LoggingMiddleware) // Ajoute le middleware de journalisation
	r.RegisterRoutes()

	log.Println("Server is running on " + fmt.Sprintf("%s:%s", commons.URL, commons.Port))
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", commons.URL, commons.Port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
