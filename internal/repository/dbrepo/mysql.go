package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/devj1003/scribble/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *mysqlDBRepo) AllUsers() bool {
	return true
}

// InsertNote inserts a note into the database
func (m *mysqlDBRepo) InsertNote(note models.Note) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO notes (user_id, title, content, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, query,
		1,
		note.Title,
		note.Content,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByID returns user by id
func (m *mysqlDBRepo) GetUserByID(id int) (models.User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, first_name, last_name, email, password, age, access_level ,created_at, updated_at
				FROM users WHERE id = ?`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.Age,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil

}

// UpdateUser updates the user data in the user table in database
func (m *mysqlDBRepo) UpdateUser(u models.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE users SET first_name=?, last_name=?, email=?, age=?, access_level=?, updated_at=?`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Age,
		u.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// Authenticate authenticates a user
func (m *mysqlDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email=?", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("Incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil

}
