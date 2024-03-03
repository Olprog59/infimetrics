package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/database"
	"github.com/Olprog59/infimetrics/models"
	"net/http"
	"regexp"
)

func SignUpEmail() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db, ok := database.FromContextDB(r)
		if !ok {
			golog.Warn("Could not get database connection from context")
		}
		err := db.Ping()
		if err != nil {
			golog.Warn("Could not ping database")
		}
		if valid, err := validateEmail(db, r.FormValue("email")); !valid {
			fmt.Fprint(w, err)
			return
		}
	}
}
func SignUpUsername() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db, ok := database.FromContextDB(r)
		if !ok {
			golog.Err("Could not get database connection")
			return
		}
		if valid, err := validateUsername(db, r.FormValue("username")); !valid {
			fmt.Fprint(w, err)
			return
		}
	}
}
func SignUpPassword() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if valid, err := validatePassword(r.FormValue("password")); !valid {
			if err == nil {
				err = errors.New("password must contain at least one uppercase letter, one lowercase letter and one number")
			}
			fmt.Fprint(w, err)
			return
		}
	}
}
func SignUpPasswordSame() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if valid, err := validateConfirmPassword(r.FormValue("password"), r.FormValue("confirm-password")); !valid {
			fmt.Fprint(w, err)
			return
		}
	}
}

// valildateEmail vérifie que l'email est valide.
func validateEmail(db *sql.DB, email string) (bool, error) {
	reg := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if reg.MatchString(email) {
		if models.IsEmailExists(db, email) {
			return false, errors.New("email already exists")
		}
		return true, nil
	}
	return false, errors.New("invalid email")
}

// validatePassword vérifie que le mot de passe est valide.
func validatePassword(password string) (bool, error) {
	if len(password) < 8 && len(password) > 64 {
		return false, errors.New("password must be between 8 and 64 characters long")
	}
	hasUpper := regexp.MustCompile(`\p{Lu}`).MatchString(password)
	hasLower := regexp.MustCompile(`\p{Ll}`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^\p{L}\p{N}]`).MatchString(password)

	if hasUpper && hasLower && hasDigit && hasSpecial {
		return true, nil
	}
	return false, errors.New("password must contain at least one uppercase letter, one lowercase letter, one number and one special character")
}

// validateConfirmPassword vérifie que le mot de passe de confirmation correspond au mot de passe.
func validateConfirmPassword(password, confirmPassword string) (bool, error) {
	if password != confirmPassword {
		return false, errors.New("passwords do not match")
	}
	return true, nil
}

// validateUsername vérifie que le nom d'utilisateur est valide.
func validateUsername(db *sql.DB, username string) (bool, error) {
	if len(username) < 3 || len(username) > 15 {
		return false, errors.New("username must be between 3 and 15 characters long")
	}

	reg := regexp.MustCompile(`^[\p{L}0-9_-]+$`)
	golog.Debug("Username: %s", username)
	if !reg.MatchString(username) {
		return false, errors.New("username must contain only letters, numbers, - and _")
	}

	if models.IsUsernameExists(db, username) {
		return false, errors.New("username already exists")
	}

	return true, nil
}
