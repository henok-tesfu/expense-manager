package models

import (
	"database/sql"

	"github.com/henok-tesfu/expense-manager/internal/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"-" validate:"required,min=6"`
}

// Register a new user
func RegisterUser(username, email, password string) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	result, err := database.DB.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		username, email, string(hashedPassword))
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Get user by email
func GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := database.DB.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// Get user by ID
func GetUserByID(id int) (*User, error) {
	user := &User{}
	err := database.DB.QueryRow("SELECT id, username, email, password FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}
