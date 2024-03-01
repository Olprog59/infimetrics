package handlers

import (
	"database/sql"
	"fmt"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/database"
	"github.com/Olprog59/infimetrics/models"
	"net/http"
)

func RegisterHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db, ok := database.FromContextDB(r)
		if !ok {
			golog.Warn("Could not get database connection from context")
		}
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				return
			}

			// Get form values
			username := r.FormValue("username")
			email := r.FormValue("email")
			password := r.FormValue("password")
			confirmPassword := r.FormValue("confirm-password")

			valid, err := ValidateRegisterForm(r, db, username, email, password, confirmPassword)
			if !valid {
				w.Header().Set("Content-Type", "text/html")
				_, err = fmt.Fprint(w, `<span id="errors" remove-content="10s">`+err.Error()+`</span>`)
				if err != nil {
					return
				}
				return
			}
			var user = new(models.UserModel)
			user.Username = username
			user.Email = email
			user.PasswordHash = password
			user.DB = &database.Db{DB: db}
			err = user.Register()
			if err != nil {
				w.Header().Set("Content-Type", "text/html")
				_, err = fmt.Fprint(w, `<span id="errors" remove-content="10s">`+err.Error()+`</span>`)
				if err != nil {
					return
				}
				return
			}
			w.Header().Set("HX-Redirect", "/sign-in")
			w.Write([]byte("User registered successfully"))

			return
		} else if r.Method == "GET" {
			renderTemplate(w, "register", &Page{
				Title: "Register",
				CSS:   []string{"sign-in-up"},
			})
			return
		}
		http.Error(w, "Not implemented", http.StatusNotImplemented)
		return
	}
}

// ValidateRegisterForm vérifie que les champs du formulaire d'inscription sont valides.
func ValidateRegisterForm(r *http.Request, db *sql.DB, username, email, password, confirmPassword string) (bool, error) {
	if valid, err := validateEmail(db, email); !valid {
		return false, err
	}
	if valid, err := validateUsername(db, password); !valid {
		return false, err
	}
	if valid, err := validatePassword(password); !valid {
		return false, err
	}
	if valid, err := validateConfirmPassword(password, confirmPassword); !valid {
		return false, err
	}
	return true, nil
}
