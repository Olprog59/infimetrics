package handlers

import (
	"fmt"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/commons"
	"github.com/Olprog59/infimetrics/database"
	"github.com/Olprog59/infimetrics/models"
	"net/http"
	"time"
)

func LoginHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db, ok := commons.FromContextDB(r)
		if !ok {
			golog.Warn("Could not get database connection from context")
		}
		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				return
			}

			// Get form values
			email := r.FormValue("email")
			password := r.FormValue("password")

			if email == "" || password == "" {
				w.Header().Set("Content-Type", "text/html")
				_, err = fmt.Fprint(w, `<span id="errors" remove-content="10s">All fields are required</span>`)
				if err != nil {
					return
				}
				return
			}

			var user = new(models.LoginModel)
			user.Email = email
			user.Password = password
			user.DB = &database.Db{DB: db}

			redis, ok := commons.FromContextRedis(r)
			if !ok {
				golog.Warn("Could not get redis connection from context")
			}
			user.Redis = redis

			ok = user.Login()
			if ok {
				// Set the cookie
				setCookieHandler(w, r, user.SessionToken)
				w.Header().Set("HX-Redirect", "/dashboard")

				w.Write([]byte("User logged in successfully"))
				return
			} else {
				w.Header().Set("Content-Type", "text/html")
				_, err = fmt.Fprint(w, `<span id="errors" remove-content="10s">Invalid email or password</span>`)
				if err != nil {
					return
				}
				return
			}
		} else if r.Method == "GET" {
			renderTemplate(w, "login", &Page{
				Title: "Login",
				CSS:   []string{"sign-in-up"},
			})
		}
	}
}

func setCookieHandler(w http.ResponseWriter, r *http.Request, value string) {
	// Création du cookie
	cookie := http.Cookie{
		Name:     "session_token",
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(commons.TimeoutCookie),
		MaxAge:   int(commons.TimeoutCookie.Seconds()),
		Secure:   true,
		HttpOnly: true,
	}

	// Ajout du cookie à la réponse
	http.SetCookie(w, &cookie)

	golog.Debug("Cookie défini")
}
