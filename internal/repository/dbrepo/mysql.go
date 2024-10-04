package dbrepo

import (
	"context"
	"time"

	"github.com/devj1003/scribble/internal/models"
)

func (m *mysqlDBRepo) AllUsers() bool {
	return true
}

// InsertNote inserts a note into the database
func (m *mysqlDBRepo) InsertNote(note models.Note) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	insertquery := `INSERT INTO notes (user_id, title, content, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, insertquery,
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
