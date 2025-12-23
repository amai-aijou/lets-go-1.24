package models

import (
	"database/sql"
	"time"
)

// 10.2: Define a User struct. Field names and types align with DB table we created
type User struct {
	ID				string
	Name			string
	Email			string
	HashedPassword	[]byte
	Created			time.Time
}

// 10.2: Define UserModel struct which wraps a database connection pool
type UserModel struct {
	DB *sql.DB
}

// 10.2: Use the Insert method to add a new record to the "users" table
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// 10.2: Use Authenticate() to verify if user exists with provided email/pass, and return relevent ID
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
