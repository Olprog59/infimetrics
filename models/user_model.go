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

func (u *UserModel) AddApplication(app ApplicationModel) {
	u.Applications = append(u.Applications, app)
}

func (u *UserModel) RemoveApplication(app ApplicationModel) {
	for i, a := range u.Applications {
		if a.AppId == app.AppId {
			u.Applications = append(u.Applications[:i], u.Applications[i+1:]...)
		}
	}
}

func (u *UserModel) GetApplication(appId uint) *ApplicationModel {
	for _, a := range u.Applications {
		if a.AppId == appId {
			return &a
		}
	}
	return nil
}

func (u *UserModel) Register() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = u.DB.DB.QueryRow("INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING user_id", u.Username, u.Email, string(hashedPassword)).Scan(&u.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (l *LoginModel) Login() bool {
	var hashedPassword string
	var userId uint
	err := l.DB.DB.QueryRow("SELECT password_hash, user_id FROM users WHERE email = $1", l.Email).Scan(&hashedPassword, &userId)
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

	golog.Debug("Session token: %s", sessionToken)
	golog.Debug("User id: %s", userId)

	err = l.Redis.SetWithTimeout(sessionToken, userId, commons.TimeoutCookie)
	if err != nil {
		log.Println("Erreur lors de la sauvegarde du token de session :", err)
		return false
	}
	return true
}

func (u *UserModel) UpdateLastLogin() error {
	_, err := u.DB.DB.Exec("UPDATE users SET last_login = $1 WHERE user_id = $2", time.Now(), u.UserId)
	if err != nil {
		return err
	}
	return nil

}

func (u *UserModel) ConvertLoginToUserModel(l *LoginModel) error {
	u.DB = l.DB
	u.Redis = l.Redis
	user_id, err := u.Redis.Get(l.SessionToken)
	if err != nil {
		return err
	}
	err = u.DB.DB.QueryRow("SELECT user_id, username, email, last_login FROM users WHERE user_id = $1", user_id).Scan(&u.UserId, &u.Username, &u.Email, &u.LastLogin)
	if err != nil {
		return err
	}
	return nil

}

func IsUsernameExists(db *sql.DB, username string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		log.Println("Erreur lors de la vérification de l'existence de l'utilisateur :", err)
		return false
	}
	return count > 0
}

func IsEmailExists(db *sql.DB, email string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
	if err != nil {
		log.Println("Erreur lors de la vérification de l'existence de l'utilisateur :", err)
		return false
	}
	golog.Debug("Email exists: %d", count)
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
