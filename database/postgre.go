package database

import (
	"database/sql"
	"github.com/olprog59/infimetrics/commons"
	"log"
)

type IDB interface {
	Connect() (*sql.DB, error)
	Close() error
}

type Db struct {
	DB sql.DB
}

func NewDB() IDB {
	return &Db{}
}

func (d *Db) Connect() (*sql.DB, error) {
	// Connect to database
	db, err := sql.Open(commons.DBDriver, commons.DBConnStr)
	if err != nil {
		log.Fatal(err)
	}

	// Check connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected!")

	_, err = db.Exec(sqlInitStr)
	if err != nil {
		return nil, err
	}
	log.Println("Successfully created table!")
	return db, nil
}

func (d *Db) Close() error {
	err := d.DB.Close()
	if err != nil {
		return err
	}
	log.Println("Database is closed")
	return nil
}

const sqlInitStr = `
CREATE TABLE IF NOT EXISTS Users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS Applications (
    app_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    app_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Logs (
    log_id SERIAL PRIMARY KEY,
    app_id INTEGER NOT NULL,
    level VARCHAR(50),
    message TEXT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB,
    FOREIGN KEY (app_id) REFERENCES Applications(app_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_users_username ON Users(username);
CREATE INDEX IF NOT EXISTS idx_applications_user_id ON Applications(user_id);
CREATE INDEX IF NOT EXISTS idx_logs_app_id ON Logs(app_id);
CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON Logs(timestamp);
`
