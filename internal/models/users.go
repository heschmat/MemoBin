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

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	q := "SELECT id, hashed_password FROM users WHERE email = ?;"

	err := m.DB.QueryRow(q, email).Scan(&id, &hashedPassword)
	// If now matchin email found, return error.
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// If password is incorrect, return error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Otherwise, the password is correct. Return the user ID.
	return id, nil
}
