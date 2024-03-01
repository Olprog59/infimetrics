package commons

import (
	"os"
	"time"
)

var (
	DB_USER     = GetEnv("DB_USER", "user")
	DB_PASSWORD = GetEnv("DB_PASSWORD", "password")
	DB_NAME     = GetEnv("DB_NAME", "dbname")
	DB_HOST     = GetEnv("DB_HOST", "localhost")
	DB_PORT     = GetEnv("DB_PORT", "5432")
	DB_DRIVER   = GetEnv("DB_DRIVER", "postgres")
	SSL_MODE    = GetEnv("SSL_MODE", "disable")

	HOST = GetEnv("HOST", "localhost")
	PORT = GetEnv("PORT", "8080")

	REDIS_HOST     = GetEnv("REDIS_HOST", "localhost")
	REDIS_PORT     = GetEnv("REDIS_PORT", "6379")
	REDIS_PASSWORD = GetEnv("REDIS_PASSWORD", "password")
)

// GetEnv is a function to get environment variable
// If the environment variable is not found, it will return the callback
func GetEnv(key, callback string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		return callback
	}
	return env
}

const TimeoutCookie = time.Hour * 24 * 7
