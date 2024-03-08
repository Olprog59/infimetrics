package api

import (
	"github.com/Olprog59/golog"
	"github.com/Olprog59/infimetrics/commons"
	"github.com/Olprog59/infimetrics/models"
	"html/template"
	"net/http"
	"time"
)

var funcs = template.FuncMap{
	"formatAsDate": formatDate,
}

func formatDate(date time.Time) string {
	return date.Format("02 Jan 2006 15:04:05")
}

func ApiAppsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		var app = new(models.ApplicationModel)
		store, ok := models.FromContextStore(r)
		if !ok {
			http.Error(w, "Internal Server Error - context db", http.StatusInternalServerError)
			return
		}
		app.Store = store

		sessionToken, err := commons.GetCookie(r, "session_token")
		if err != nil {
			http.Error(w, "Internal Server Error - context session_token", http.StatusInternalServerError)
			return
		}
		session, err := app.Store.HGet(sessionToken, "userID")
		if err != nil {
			http.Error(w, "Internal Server Error - redis get userID", http.StatusInternalServerError)
			return
		}
		app.UserId = commons.StringToUint(session)

		apps, err := app.FindAllApps()
		if err != nil {
			golog.Err(err.Error())
			http.Error(w, "Internal Server Error - find all apps", http.StatusInternalServerError)
			return
		}
		//golog.Info("Apps found: %+v", apps)
		if len(apps) == 0 {
			return
		}

		t := template.Must(template.New("apps").Funcs(funcs).Parse(templateApiApps))
		err = t.Execute(w, struct {
			Data []models.ApplicationModel
		}{
			Data: apps,
		})
		if err != nil {
			golog.Err(err.Error())
			http.Error(w, "Internal Server Error - write apps found", http.StatusInternalServerError)
			return
		}

	}
}

const templateApiApps = `
{{ range .Data }}
	<tr class="app" class="fade-in">
		<td data-label="AppName">{{ .AppName }}</td>
		<td data-label="Description">{{ .Description }}</td>
		<td data-label="CreatedAt">{{ .CreatedAt | formatAsDate }}</td>
		<td data-label="Actions">
			<a hx-get="/api/v1/apps/{{.Token}}" hx-target="body" hx-swap="outerHTML" hx-push-url="true" class="view">View</a>
			<a hx-delete="/api/v1/app/{{.Token}}" class="delete">Delete</a>
		</td>
	</tr>
{{ end }}
`

func ApiNewAppsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		t := template.Must(template.New("apps").Parse(templateApiNewApps))
		err := t.Execute(w, nil)
		if err != nil {
			golog.Err(err.Error())
			http.Error(w, "Internal Server Error - write apps found", http.StatusInternalServerError)
			return
		}
	}
}

const templateApiNewApps = `
<form hx-post="/api/v1/app" hx-swap="#form-errors">
	<input type="text" name="app_name"
			placeholder="App name"
			minlength="3"
			maxlength="50"

			required>
	<textarea name="description"
			placeholder="Description"
			minlength="10"
			maxlength="200"
			pattern="^\w{3,200}$"
			spellcheck="true"
			autocomplete="off"
			required></textarea>
	<div id="form-errors"></div>
	<div>
		<button type="submit" class="create">Create</button>
		<button type="button" hx-delete="/api/v1/app/modal" hx-target="#apps__new" hx-swap="innerHTML swap:1s" class="close">Close</button>
	</div>
</form>
`

func ApiNewAppsPostHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		r.ParseForm()

		var app = new(models.ApplicationModel)
		store, ok := models.FromContextStore(r)
		if !ok {
			http.Error(w, "Internal Server Error - context db", http.StatusInternalServerError)
			return
		}

		app.Store = store
		sessionToken, err := commons.GetCookie(r, "session_token")
		if err != nil {
			http.Error(w, "Internal Server Error - context session_token", http.StatusInternalServerError)
			return
		}

		session, err := app.Store.HGet(sessionToken, "userID")
		if err != nil {
			http.Error(w, "Internal Server Error - redis get userID", http.StatusInternalServerError)
			return
		}

		app.UserId = commons.StringToUint(session)
		app.AppName = r.FormValue("app_name")
		app.Description = r.FormValue("description")

		err = app.Verification()
		if err != nil {
			golog.Err(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = app.CreateApp()
		if err != nil {
			golog.Err(err.Error())
			http.Error(w, "Internal Server Error - create app", http.StatusInternalServerError)
			return
		}
		//w.Header().Set("HX-Trigger", "newApp")
		//_, err = fmt.Fprintf(w, "App created: %s", app.AppName)
		//if err != nil {
		//	return
		//}
		app.CreatedAt = time.Now()
		t := template.Must(template.New("apps").Funcs(funcs).Parse(templateApiNewApp))
		err = t.Execute(w, app)
		if err != nil {
			golog.Err(err.Error())
			http.Error(w, "Internal Server Error - write apps found", http.StatusInternalServerError)
			return
		}
		err = app.InsertAppMongo()
		if err != nil {
			return
		}
	}
}

const templateApiNewApp = `
<tr class="app" hx-swap-oob="beforeend:#apps__all">
	<td data-label="AppName">{{ .AppName }}</td>
	<td data-label="Description">{{ .Description }}</td>
	<td data-label="CreatedAt">{{ .CreatedAt | formatAsDate }}</td>
	<td data-label="Actions">
		<a hx-get="/api/v1/apps/{{.Token}}" hx-target="body" hx-swap="outerHTML" hx-push-url="true" class="view">View</a>
		<a hx-delete="/api/v1/app/{{.Token}}" class="delete">Delete</a>
	</td>
</tr>
`

func ApiNewAppsDeleteHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		golog.Info("Delete app")
		w.Write([]byte(""))
	}
}

func ApiDeleteAppsPostHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		store, ok := models.FromContextStore(r)
		if !ok {
			http.Error(w, "Internal Server Error - context db", http.StatusInternalServerError)
			return
		}
		token := r.PathValue("token")
		app := new(models.ApplicationModel)
		app.Store = store
		app.Token = token
		err := app.DeleteAppMongo()
		if err != nil {
			http.Error(w, "Internal Server Error - delete app", http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Redirect", "/apps")
		w.Write([]byte(""))
	}
}
