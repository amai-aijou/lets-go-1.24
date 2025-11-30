package models

import (
	"database/sql"
	"errors"
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
	//the backquote (`) character is used to split a single line into multiple lines in the code
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
	// 4.7: Execute SQL statement to pull & read contents of table
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Runs QueryRow() method against connection pool to execute SQL statement. Passes untrusted id variable
	// as the value for the ? placeholders send previously (prevents SQL injection attacks). Returns point to a sql.Row value
	// which holds the result from the database
	row := m.DB.QueryRow(stmt, id)

	// Initialize a new "zeroed" Snippet struct.
	var s Snippet
	
	// row.Scan() copies values from fields in sql.Row to the Snippet struct in main.go. must be same args returned by statement
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// If the query returns nothing, row.Scan() will return a sql.ErrNoRows error. errors.Is() returns our own ErrNoRecord error
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil

}

// Return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}

