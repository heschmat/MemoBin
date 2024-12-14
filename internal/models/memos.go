package models

import (
	"database/sql"
	"time"
)

// Define a `Memo` type to hold the data or an individual "memo".
type Memo struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// `MemoModel` wraps a sql.DB connection pool.
type MemoModel struct {
	DB      *sql.DB
}

// Insert
func (m *MemoModel) Insert(title string, content string, expires int) (int, error) {
	// Using `` we can split the query we want to execute over multiple lines for readability.
	// N.B. PostgreSQL uses $N notation for placeholder parameter.
	query := `INSERT INTO memos (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY));`

	// `Exec()` returns a sql.Result type
	// This contains basic information about what happened when the query executed.
	result, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	// Get the ID of newly inserted record in the *memos* table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
