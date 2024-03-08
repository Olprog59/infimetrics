package handlers

import (
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/commons"
	"github.com/Olprog59/infimetrics/models"
	"net/http"
)

func LoginHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, ok := models.FromContextStore(r)
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
				//w.Header().Set("Content-Type", "text/html")
				//_, err = fmt.Fprint(w, `<span id="errors" remove-content="10s">All fields are required</span>`)
				//if err != nil {
				//	return
				//}
				http.Error(w, "All fields are required", http.StatusInternalServerError)
				return
			}

			var login = new(models.LoginModel)
			login.Email = email
			login.Password = password
			login.Store = store

			ok = login.Login()
			if ok {
				// Set the cookie
				var user = new(models.UserModel)
				err := user.ConvertLoginToUserModel(login)
				if err != nil {
					golog.Err("Could not convert login to user model")
					return
				}
				//setCookieHandler(w, r, Username, user.Username)
				err = user.UpdateLastLogin()
				if err != nil {
					golog.Err("Could not update last login")
					return
				}
				commons.SetCookie(w, r, commons.SessionToken, login.SessionToken)
				w.Header().Set("HX-Redirect", "/")
				return
			} else {
				//w.Header().Set("Content-Type", "text/html")
				//_, err = fmt.Fprint(w, `<span id="errors" remove-content="10s">Invalid email or password</span>`)
				//if err != nil {
				//	return
				//}
				http.Error(w, "Invalid email or password", http.StatusInternalServerError)
				return
			}
		} else if r.Method == "GET" {
			commons.RenderTemplate(w, r, "login", &commons.Page{
				Title: "Login",
			})
		}
	}
}
