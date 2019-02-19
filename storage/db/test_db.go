package db

import (
	"github.com/jinzhu/gorm"
)

type TestDB struct{}

func (i TestDB) Open() *gorm.DB {
	db, err := gorm.Open("sqlite3", "/tmp/test.db")
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
