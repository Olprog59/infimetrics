package commons

import "os"

var (
	// DBHost is a variable to store database host
	DBHost = GetEnv("POSTGRES_HOST", "localhost")
	// DBUser is a variable to store database user
	DBUser = GetEnv("POSTGRES_USER", "postgres")
	// DBPassword is a variable to store database password
	DBPassword = GetEnv("POSTGRES_PASSWORD", "postgres")
	// DBName is a variable to store database name
	DBName = GetEnv("POSTGRES_DB", "todos")
	// DBPort is a variable to store database port
	DBPort = GetEnv("POSTGRES_PORT", "5432")
	// DBSSLMode is a variable to store database ssl mode
	DBSSLMode = GetEnv("DB_SSL_MODE", "disable")

	// DBConnStr is a variable to store database connection string
	DBConnStr = "postgresql://" + DBUser + ":" + DBPassword + "@" + DBHost + "/" + DBName + "?sslmode=" + DBSSLMode

	// DBDriver is a variable to store database driver
	DBDriver = "postgres"

	// URL http server
	URL = GetEnv("URL", "")

	// Port is a variable to store server port
	Port = GetEnv("PORT", "8080")

	// SecretKey is a variable to store secret key
	SecretKey = GetEnv("secret_key", "secret")

	// TokenDuration is a variable to store token duration
	TokenDuration = GetEnv("token_duration", "1h")

	// TokenRefreshDuration is a variable to store token refresh duration
	TokenRefreshDuration = GetEnv("token_refresh", "24h")
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
