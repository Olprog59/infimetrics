package handlers

import (
	"github.com/Olprog59/infimetrics/database"
	"net/http"
)

func WithRedisHandler(redis *database.RedisDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := redis.Ping()
		if err != nil {
			http.Error(w, "Error connecting to Redis", http.StatusInternalServerError)
			return
		}

		err = redis.Set("key", "value")
		if err != nil {
			http.Error(w, "Error setting value in Redis", http.StatusInternalServerError)
			return
		}

		val, err := redis.Get("key")
		if err != nil {
			http.Error(w, "Error getting value from Redis", http.StatusInternalServerError)
			return
		}

		_, err = w.Write([]byte(val))
		if err != nil {
			return
		}
	}
}
