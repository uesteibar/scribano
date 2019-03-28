package db

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres driver
	"github.com/xo/dburl"
)

// Database defines the interface of a db client
type Database interface {
	Open() *gorm.DB
}

// DB is the db interface
type DB struct{}

// Open opens a connection to the database
func (i DB) Open() *gorm.DB {
	url := os.Getenv("PG_URL")
	u, err := dburl.Parse(url)
	db, err := gorm.Open("postgres", u.DSN)
	if err != nil {
		panic(err)
	}

	return db
}
