package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	// Create a *bcrypt* hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	q := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(q, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		// N.B. The *mySQLError* variable is initialized to *nil*.
		// It gets populated ONLY IF `errors.As` detects that *err* is of type `*mysql.MySQLError`.
		// If `errors.As` succeeds, it assigns the underlying `*mysql.MySQLError` object to *mySQLError*,
		// enabling the function to access its `Number` & `Message` fields for detailed error handling.
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 162 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}
