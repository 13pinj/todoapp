package models

import (
	"github.com/13pinj/todoapp/core/log"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var DB gorm.DB

func init() {
	db, err := gorm.Open("sqlite3", "debug.db")
	if err != nil {
		panic(err)
	}
	DB = db
	DB.LogMode(true)
	DB.SetLogger(gorm.Logger{LogWriter: log.Logger()})
}
