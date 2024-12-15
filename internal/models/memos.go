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

func (m *MemoModel) Latest() ([]Memo, error) {
	query := `SELECT id, title, content, created, expires FROM memos
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10;`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	// *defer* rows.Close() to ensure the sql.Rows resultset is always properly closed
	// before Latest() method returns.
	// N.B. This should come **after** checking for an error from the Query() method.
	// Or you may get a **panic**.
	defer rows.Close()

	// Initialize an empty slice to hold the Memo structs.
	var memos []Memo

	for rows.Next() {
		var m Memo
		err = rows.Scan(&m.ID, &m.Title, &m.Content, &m.Created, &m.Expires)
		if err != nil {
			return nil, err
		}

		memos = append(memos, m)
	}

	// Retrieve any error that was encountered dureing the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went OK, return the Memos slice.
	return memos, nil
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
