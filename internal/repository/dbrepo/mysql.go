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

// InsertNewUser inserts a new user into the database
func (m *mysqlDBRepo) InsertNewUser(user models.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO users (first_name, last_name, email, password, access_level, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		1,
		time.Now(),
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

// InsertNote inserts a note into the database
func (m *mysqlDBRepo) InsertNote(note models.Note, id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO notes (user_id, title, content, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, query,
		id,
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

// ViewNoteAtIndex returns all notes by id for index page
func (m *mysqlDBRepo) ViewNoteAtIndex(id int) ([]models.Note, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, title, content, created_at, updated_at FROM notes
				WHERE user_id = ?`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var viewnoteatindex []models.Note

	// iterate over the rows
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err // Handle scanning error
		}
		viewnoteatindex = append(viewnoteatindex, note)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return viewnoteatindex, nil

}

// ViewNoteData shows the note data using note id
func (m *mysqlDBRepo) ViewNoteData(nid int, id int) (models.Note, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT * FROM notes WHERE id=? AND user_id=?`

	row := m.DB.QueryRowContext(ctx, query, nid, id)
	var notedata models.Note
	err := row.Scan(
		&notedata.ID,
		&notedata.UserID,
		&notedata.Title,
		&notedata.Content,
		&notedata.CreatedAt,
		&notedata.UpdatedAt,
	)
	if err != nil {
		return notedata, err
	}

	return notedata, nil

}

// DeleteNote deletes a note from the database
func (m *mysqlDBRepo) DeleteNote(nid int, id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM notes WHERE id=? AND user_id=?`
	_, err := m.DB.ExecContext(ctx, query, nid, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateNote updates a note in the database
func (m *mysqlDBRepo) UpdateNote(u models.Note, nid int, id int) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE notes SET title=?, content=?, updated_at=? WHERE id=? AND user_id=?`
	_, err := m.DB.ExecContext(ctx, query,
		u.Title,
		u.Content,
		time.Now(),
		nid,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}
