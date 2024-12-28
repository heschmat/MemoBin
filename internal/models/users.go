package models

import (
	"database/sql"
	"time"
)

// Field names & types match with those of *users* table:
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// This wraps a database connection pool.
type UserModel struct {
	DB *sql.DB
}

// Add a new record to the *users* table.
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}
