package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type TestDB struct {
	name string
}

func GetUniqueDB() TestDB {
	name := uuid.New().String()
	return TestDB{name: name}
}

func (i TestDB) Open() *gorm.DB {
	dbPath := fmt.Sprintf("/tmp/%s_test.db", i.name)
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
