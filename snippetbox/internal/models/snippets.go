package models

import (
	"database/sql"
	"time"
)

// Define a Snippet type to hold the data for individual snippets
// Fields correspond to the MySQL snippets table
type Snippet struct {
	ID		int
	Title	string
	Content string
	Created time.Time
	Expires time.Time
}

// Define SnippetModel type to wrap a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Add SQL statement/s to execute:
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use Exec() on the connectino pool to execute the statement.
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// User LastInsertID() method to get the ID of the record
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// return a specific snippet based on its id
func (m *SnippetModel) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

// Return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}

