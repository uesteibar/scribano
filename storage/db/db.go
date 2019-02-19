package db

import (
	"github.com/jinzhu/gorm"
)

type Database interface {
	Open() *gorm.DB
}

type DB struct{}

func (i DB) Open() *gorm.DB {
	db, err := gorm.Open("sqlite3", "/tmp/asyncapi.db")
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
