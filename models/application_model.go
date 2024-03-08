package models

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Olprog59/golog"
	"github.com/google/uuid"
	"log"
	"regexp"
	"time"
)

type ApplicationModel struct {
	Store       *Store    `bson:"-"`
	AppId       uint      `bson:"appId"`
	UserId      uint      `bson:"userId"`
	AppName     string    `bson:"appName"`
	Token       string    `bson:"token"`
	Description string    `bson:"description"`
	CreatedAt   time.Time `bson:"createdAt"`
}

const (
	// findAllApps récupère toutes les applications d'un utilisateur en postgresql
	findAllApps = `SELECT app_id, user_id, app_name, token, created_at, description FROM application WHERE user_id = $1`
	// createApp crée une nouvelle application en postgresql
	createApp = `INSERT INTO application (user_id, app_name, token, description) VALUES ($1, $2, $3, $4)`

	getAppByToken = `SELECT app_id, user_id, app_name, token, description, created_at FROM application WHERE token = $1`

	deleteApp = `DELETE FROM application WHERE token = $1`
)

func (app *ApplicationModel) FindAllApps() ([]ApplicationModel, error) {
	stmt, err := app.Store.DB.Prepare(findAllApps)
	if err != nil {
		return nil, err
	}

	defer (func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("Erreur lors de la fermeture du statement :", err)
		}
	})(stmt)

	rows, err := stmt.Query(app.UserId)
	if err != nil {
		return nil, err
	}

	var apps []ApplicationModel
	for rows.Next() {
		var a ApplicationModel
		err = rows.Scan(&a.AppId, &a.UserId, &a.AppName, &a.Token, &a.CreatedAt, &a.Description)
		if err != nil {
			return nil, err
		}
		apps = append(apps, a)
	}
	return apps, nil
}

func (app *ApplicationModel) CreateApp() error {
	stmt, err := app.Store.DB.Prepare(createApp)
	if err != nil {
		return err
	}

	defer (func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("Erreur lors de la fermeture du statement :", err)
		}
	})(stmt)

	app.Token = uuid.NewString()

	_, err = stmt.Exec(app.UserId, app.AppName, app.Token, app.Description)
	if err != nil {
		return err
	}
	return nil
}

func (app *ApplicationModel) Verification() error {
	if len(app.AppName) < 3 || len(app.AppName) > 50 {
		return errors.New("App name invalid")
	}
	reg := regexp.MustCompile(`^\p{L}[\p{L}-_0-9]{1,48}\p{L}$`)
	if !reg.MatchString(app.AppName) {
		return errors.New("App name invalid")
	}

	if len(app.Description) < 10 || len(app.Description) > 200 {
		return errors.New("Description invalid")
	}
	return nil
}

func (app *ApplicationModel) GetUserIdByToken() (string, error) {
	var userId string
	stmt, err := app.Store.DB.Prepare(getAppByToken)
	if err != nil {
		return "", err
	}
	// app_id, user_id, app_name, token, description, created_at
	err = stmt.QueryRow(app.Token).Scan(&app.AppId, &userId, &app.AppName, &app.Token, &app.Description, &app.CreatedAt)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (app *ApplicationModel) InsertAppMongo() error {
	client, err := app.Store.Mongo.Connect()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.Database(app.Token).Collection("apps").InsertOne(ctx, app)
	if err != nil {
		return err
	}
	golog.Info("Inserted a single document: %v", res.InsertedID)
	return nil
}

func (app *ApplicationModel) DeleteAppMongo() error {
	client, err := app.Store.Mongo.Connect()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// remove collection
	err = client.Database(app.Token).Collection("apps").Drop(ctx)
	if err != nil {
		return err
	}

	stmt, err := app.Store.DB.Prepare(deleteApp)
	if err != nil {
		return err
	}

	defer (func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("Erreur lors de la fermeture du statement :", err)
		}
	})(stmt)

	_, err = stmt.Exec(app.Token)
	if err != nil {
		return err
	}

	golog.Info("Deleted database")
	return nil
}
