package models

import (
	"github.com/Olprog59/infimetrics/database"
)

type Store struct {
	*database.Db
	*database.RedisDB
	*database.Mongo
}
