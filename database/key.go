package database

// Context keys
type contextKey string

const (
	DbKey    = contextKey("db")
	RedisKey = contextKey("redis")
)
