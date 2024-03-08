package database

import (
	"database/sql"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/commons"
)

type IDB interface {
	Connect() (*sql.DB, error)
	Close() error
}

type Db struct {
	DB *sql.DB
}

func NewDB() *Db {
	return &Db{}
}

func (d *Db) Connect() (*Db, error) {
	dbConnStr := "postgresql://" + commons.DB_USER + ":" + commons.DB_PASSWORD + "@" + commons.DB_HOST + "/" + commons.DB_NAME + "?sslmode=" + commons.SSL_MODE
	// Connect to database
	db, err := sql.Open(commons.DB_DRIVER, dbConnStr)
	if err != nil {
		golog.Err(err.Error())
		return nil, err
	}

	// Check connection
	err = db.Ping()
	if err != nil {
		golog.Err(err.Error())
		return nil, err
	}
	golog.Success("Successfully connected!")

	_, err = db.Exec(sqlInitStr)
	if err != nil {
		return nil, err
	}
	golog.Success("Successfully created table!")
	var dbInstance = new(Db)
	dbInstance.DB = db
	return dbInstance, nil
}

func (d *Db) Close() error {
	err := d.DB.Close()
	if err != nil {
		return err
	}
	golog.Err("Database is closed")
	return nil
}

const sqlInitStr = `
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS application (
    app_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    app_name VARCHAR(100) NOT NULL,
    token VARCHAR(64) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_applications_user_id ON application(user_id);
`
