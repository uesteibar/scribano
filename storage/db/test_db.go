package db

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// TestDB is a test-only database interface
type TestDB struct {
	name string
}

// GetUniqueDB returns a unique database client for isolated tests
func GetUniqueDB() TestDB {
	name := uuid.New().String()
	return TestDB{name: name}
}

// Open a connection to the database
func (i TestDB) Open() *gorm.DB {
	dbPath := fmt.Sprintf("/tmp/%s_test.db", i.name)
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
