package database

import (
	"database/sql"
	"net/http"
)

func FromContextDB(r *http.Request) (*sql.DB, bool) {
	db, ok := r.Context().Value(DbKey).(*sql.DB)
	return db, ok
}

func FromContextRedis(r *http.Request) (*RedisDB, bool) {
	redis, ok := r.Context().Value(RedisKey).(*RedisDB)
	return redis, ok
}
