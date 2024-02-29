package commons

import (
	"database/sql"
	"github.com/Olprog59/infimetrics/database"
	"net/http"
)

func FromContextDB(r *http.Request) (*sql.DB, bool) {
	db, ok := r.Context().Value(database.DbKey).(*sql.DB)
	return db, ok
}

func FromContextRedis(r *http.Request) (*database.RedisDB, bool) {
	redis, ok := r.Context().Value(database.RedisKey).(*database.RedisDB)
	return redis, ok
}
