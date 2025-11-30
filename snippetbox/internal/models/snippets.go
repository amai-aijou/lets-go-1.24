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

	// 4.7 QueryRow() is used for queries that return a singular row. Passes untrusted id variable
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
	// 4.8: SQL Statement to pull 10 latest Snippets (ordered by ID)
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// 4.8: DB.query() is used for queries that return multiple rows. returns sql.Rows resultset with the rows requested
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// 4.8: Defer rows.Close() to ensure the sql.rows resultset is properly closed before method Latest() returns.
	// Defer statements come *after* error-checking from Query(); otherwise, trying to close a nil resultset causes a panic.
	defer rows.Close()

	// 4.8: Initialize an empty slice to hold the Snippet structs
	var snippets []Snippet

	// 4.8: rows.Next iterates through the rows in the resultset, and prepare each row to be acted on by rows.Scan().
	// If iteration completes, the resultset auto-closes itself and frees up the underlying database connection.
	for rows.Next() {
		// 4.8: create a new zero value Snippet struct.
		var s Snippet
		// 4.8: Use rows.Scan() to copy values from each field in the row to the new struct we created
		// Args are pointers to the place we want to copy data into, and number of args much match exactly
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// 4.8: Append to the slice of snippets
		snippets = append(snippets, s)
	}

	// 4.8: When rows.Next() loop completes, call rows.Err() to retrieve any errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// 4.8: If all was successful, return the Snippets slice!
	return snippets, nil
}

