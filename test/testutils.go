package testutils

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) *gorm.DB {
	dsn := "host=localhost user=postgres password=password dbname=library port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error. Failed to connect to test DB: %v", err)
	}

	schema := "test_" + t.Name()
	db.Exec("drop schema if exists " + schema + " cascade")
	db.Exec("create schema " + schema)
	db.Exec("set search_path to " + schema)

	return db
}

func FreeTestDB(t *testing.T, db *gorm.DB) {
	schema := "test_" + t.Name()
	db.Exec("drop schema if exists " + schema + " cascade")
}
