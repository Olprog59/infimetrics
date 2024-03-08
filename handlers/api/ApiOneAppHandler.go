package api

import (
	"errors"
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/commons"
	"github.com/Olprog59/infimetrics/models"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func ApiOneAppHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//if r.Header.Get("HX-Request") != "true" {
		//	ref := r.Header.Get("referer")
		//	golog.Info("Ref: %s", ref)
		//	http.Redirect(w, r, ref, http.StatusSeeOther)
		//	return
		//}
		sessionToken, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Internal Server Error - context session_token", http.StatusInternalServerError)
			return
		}
		token := r.PathValue("token")
		w.Header().Set("Content-Type", "text/html")
		store, ok := models.FromContextStore(r)
		if !ok {
			http.Error(w, "Internal Server Error - context db", http.StatusInternalServerError)
			return
		}
		userID, err := store.HGet(sessionToken.Value, "userID")
		if err != nil {
			http.Error(w, "Internal Server Error - redis get userID", http.StatusInternalServerError)
			return
		}
		app := new(models.ApplicationModel)
		app.Store = store
		app.Token = token
		idByToken, err := app.GetUserIdByToken()
		if err != nil {
			return
		}
		if idByToken != userID {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log := &models.LogModel{
			Store: store,
		}

		logs, err := log.GetLogsByAppToken(app.Token)

		if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
			golog.Err("Error getting logs: %s", err)
			return
		}

		page := &commons.Page{
			Title:    "One app",
			Data:     logs,
			AppToken: token,
		}
		commons.RenderTemplate(w, r, "oneApp", page)
		return
	}
}
