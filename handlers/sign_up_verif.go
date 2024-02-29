package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/commons"
	"github.com/Olprog59/infimetrics/models"
	"net/http"
	"regexp"
)

func SignUpEmail() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db, ok := commons.FromContextDB(r)
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
		db, ok := commons.FromContextDB(r)
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
	if len(password) < 8 {
		return false, errors.New("password must be at least 8 characters long")
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)

	return hasUpper && hasLower && hasDigit, nil
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
	if len(username) < 3 {
		return false, errors.New("username must be at least 3 characters long")
	}

	if models.IsUsernameExists(db, username) {
		return false, errors.New("username already exists")
	}

	return true, nil
}
