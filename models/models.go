package models

import (
	"os"

	"github.com/13pinj/todoapp/core/log"

	"github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/jinzhu/gorm"
	// SQLite driver for gorm
	_ "github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/mattn/go-sqlite3"
	// PostreSQL driver for gorm
	_ "github.com/13pinj/todoapp/Godeps/_workspace/src/github.com/lib/pq"
)

// DB - переменная для доступа к БД
var DB gorm.DB

func init() {
	mode := os.Getenv("TODOAPP_MODE")
	var err error
	if mode == "HEROKU" {
		DB, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	} else {
		DB, err = gorm.Open("sqlite3", "debug.db")
	}
	if err != nil {
		panic(err)
	}
	DB.LogMode(true)
	DB.SetLogger(gorm.Logger{LogWriter: log.Logger()})
}
