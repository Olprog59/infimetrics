package api

import (
	"github.com/Olprog59/infimetrics/models"
	"net/http"
	"time"
)

func ApiLogsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.PathValue("token")
		if token == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Get the store from the request context
		store, ok := models.FromContextStore(r)
		if !ok {
			http.Error(w, "Internal Server Error - context db", http.StatusInternalServerError)
			return
		}

		log := models.NewLogModel(store, "INFO", "Test", time.Now(), "Test")
		err := log.InsertLogMongo(token)
		if err != nil {
			http.Error(w, "Internal Server Error - insert log", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

//func ApiLogsWatchHandler() func(http.ResponseWriter, *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		token := r.PathValue("token")
//		if token == "" {
//			http.Error(w, "Bad Request", http.StatusBadRequest)
//			return
//		}
//
//		// Get the store from the request context
//		store, ok := models.FromContextStore(r)
//		if !ok {
//			http.Error(w, "Internal Server Error - context db", http.StatusInternalServerError)
//			return
//		}
//
//		log := new(models.LogModel)
//		log.Store = store
//		coll := log.GetConnection(token)
//		pipeline := mongo.Pipeline{bson.D{{"$match", bson.D{{"operationType", "insert"}}}}}
//		changeStream, err := coll.Watch(r.Context(), pipeline)
//		if err != nil {
//			golog.Err("Error watching logs: %s", err)
//			http.Error(w, "Internal Server Error - watch", http.StatusInternalServerError)
//			return
//		}
//
//		defer changeStream.Close(r.Context())
//
//		for changeStream.Next(r.Context()) {
//			var data models.LogModel
//			if err := changeStream.Decode(&data); err != nil {
//				http.Error(w, "Internal Server Error - decode", http.StatusInternalServerError)
//				return
//			}
//			commons.RenderTemplate(w, r, "log", &commons.Page{
//				Title: "Log",
//				Data:  data,
//			})
//		}
//	}
//}
