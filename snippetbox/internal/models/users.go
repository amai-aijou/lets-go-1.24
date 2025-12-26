package models

import (
	"database/sql"
	"errors" // 10.3
	"strings" // 10.3
	"time"

	"github.com/go-sql-driver/mysql" // 10.3
	"golang.org/x/crypto/bcrypt" //10.3: Used for hashing and decrypting passwords

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
	// 10.3: Create a bcrypt hash of plain-text password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	
	// 10.3: Creates a Prepared Statement for SQL, omitting the potentially-sensitive name, email, and hashed password values for later insertion
	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`
	
	// 10.3: Use Exec() method to insert the user details and hashed password into the users table
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// 10.3: If error exists and is type *mysql.MySQLError, assign to mySQLError variable, then check for a 1062 error code; if found, return ErrDuplicateEmail
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

// 10.2: Use Authenticate() to verify if user exists with provided email/pass, and return relevent ID
func (m *UserModel) Authenticate(email, password string) (int, error) {
	// 10.4: Retrieve Id and hashed password for given email. If no email, return ErrInvalidCredentials error
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// 10.4: Check whether hashed password and plain-text password match. If not, return ErrInvalidCredentials error
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}else {
			return 0, err
		}
	}

	// 10.4: If password is correct, return the user ID
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
