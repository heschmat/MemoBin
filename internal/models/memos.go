package models

import (
	"database/sql"
	"errors"
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

// GET memo/{id}
func (m *MemoModel) Get(id int) (Memo, error) {
	query := `SELECT id, title, content, created, expires FROM memos
	WHERE expires > UTC_TIMESTAMP() AND id = ?;`

	// Returns a pointer to a `sql.Row` object, which holds the result.
	row := m.DB.QueryRow(query, id)

	var memo Memo // Initialize a new zeroed Memo struct.

	err := row.Scan(&memo.ID, &memo.Title, &memo.Content, &memo.Created, &memo.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Memo{}, ErrNoRecord // We define `ErrNoRecord`
		}
		return Memo{}, err
	}

	// If everythig went OK, return the filled Memo struct.
	return memo, nil
}

// POST
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
