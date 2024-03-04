package models

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/commons"
	"github.com/Olprog59/infimetrics/database"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

const (
	insertUserQuery          = "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING user_id"
	selectUserQuery          = "SELECT password_hash, user_id, username FROM users WHERE email = $1"
	updateUserQuery          = "UPDATE users SET last_login = $1 WHERE user_id = $2"
	selectUserByIdQuery      = "SELECT user_id, username, email, last_login FROM users WHERE user_id = $1"
	countUserByUsernameQuery = "SELECT COUNT(*) FROM users WHERE username = $1"
	countUserByEmailQuery    = "SELECT COUNT(*) FROM users WHERE email = $1"
)

type UserModel struct {
	DB           *database.Db       `json:"-"`
	Redis        *database.RedisDB  `json:"-"`
	UserId       uint               `json:"user_id"`
	Username     string             `json:"username"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"passwordHash"`
	CreatedAt    time.Time          `json:"createdAt"`
	LastLogin    *time.Time         `json:"lastLogin"`
	Applications []ApplicationModel `json:"applications"`
}

type LoginModel struct {
	DB           *database.Db      `json:"-"`
	Redis        *database.RedisDB `json:"-"`
	Email        string            `json:"email"`
	Password     string            `json:"passwordHash"`
	SessionToken string            `json:"-"`
}

// ... (rest of the code remains the same)

func (u *UserModel) Register() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt, err := u.DB.DB.Prepare(insertUserQuery)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("Erreur lors de la fermeture du statement :", err)
		}
	}(stmt)

	err = stmt.QueryRow(u.Username, u.Email, string(hashedPassword)).Scan(&u.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (l *LoginModel) Login() bool {
	var hashedPassword string
	var userId uint
	var username string

	stmt, err := l.DB.DB.Prepare(selectUserQuery)
	if err != nil {
		log.Println("Erreur lors de la préparation de la requête :", err)
		return false
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("Erreur lors de la fermeture du statement :", err)
		}
	}(stmt)

	err = stmt.QueryRow(l.Email).Scan(&hashedPassword, &userId, &username)
	if err != nil {
		log.Println("Erreur lors de la récupération de l'utilisateur :", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(l.Password))
	if err != nil {
		golog.Err("Erreur lors de la comparaison des mots de passe : %s", err.Error())
		return false
	}

	sessionToken := generateSessionToken()
	l.SessionToken = sessionToken
	err = l.Redis.HSetWithTimeout(sessionToken, map[string]any{"userID": userId, "username": username}, commons.TimeoutCookie)
	//err = l.Redis.SetWithTimeout(sessionToken, userId, commons.TimeoutCookie)
	if err != nil {
		log.Println("Erreur lors de la sauvegarde du token de session :", err)
		return false
	}
	return true
}

func (u *UserModel) UpdateLastLogin() error {
	stmt, err := u.DB.DB.Prepare(updateUserQuery)
	if err != nil {
		log.Println("Erreur lors de la préparation de la requête :", err)
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("Erreur lors de la fermeture du statement :", err)
		}
	}(stmt)
	row := stmt.QueryRow(time.Now(), u.UserId)
	if row.Err() != nil {
		return err
	}
	return nil

}

func (u *UserModel) ConvertLoginToUserModel(l *LoginModel) error {
	u.DB = l.DB
	u.Redis = l.Redis
	userId, err := u.Redis.HGet(l.SessionToken, "userID")
	if err != nil {
		return err
	}
	stmt, err := u.DB.DB.Prepare(selectUserByIdQuery)
	if err != nil {
		log.Println("Erreur lors de la préparation de la requête :", err)
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("Erreur lors de la fermeture du statement :", err)
		}
	}(stmt)
	err = stmt.QueryRow(userId).Scan(&u.UserId, &u.Username, &u.Email, &u.LastLogin)
	if err != nil {
		return err
	}
	return nil

}

func IsUsernameExists(db *sql.DB, username string) bool {
	var count int
	stmt, err := db.Prepare(countUserByUsernameQuery)
	if err != nil {
		log.Println("Erreur lors de la préparation de la requête :", err)
		return false
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("Erreur lors de la fermeture du statement :", err)
		}
	}(stmt)
	err = stmt.QueryRow(username).Scan(&count)
	if err != nil {
		log.Println("Erreur lors de la vérification de l'existence de l'utilisateur :", err)
		return false
	}
	return count > 0
}

func IsEmailExists(db *sql.DB, email string) bool {
	var count int
	stmt, err := db.Prepare(countUserByEmailQuery)
	if err != nil {
		log.Println("Erreur lors de la préparation de la requête :", err)
		return false
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("Erreur lors de la fermeture du statement :", err)
		}
	}(stmt)
	err = stmt.QueryRow(email).Scan(&count)
	if err != nil {
		log.Println("Erreur lors de la vérification de l'existence de l'utilisateur :", err)
		return false
	}
	return count > 0
}

func generateSessionToken() string {
	// Définit la taille du token. 32 octets donnent 64 caractères hexadécimaux.
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Erreur lors de la génération du token de session : %v", err)
	}

	// Retourne le token sous forme de chaîne hexadécimale.
	return hex.EncodeToString(b)
}
