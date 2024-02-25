package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func WithDBHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, ok := r.Context().Value("db").(*sql.DB)
		if !ok {
			http.Error(w, "Could not get database connection from context", http.StatusInternalServerError)
			return
		}

		result, err := db.Exec("SELECT 1")
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = result.RowsAffected()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println(result)

		// Ici, tu peux utiliser db pour interagir avec ta base de données
		fmt.Fprintf(w, "Database connection retrieved successfully")
	}
}
