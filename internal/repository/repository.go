package repository

import "github.com/devj1003/scribble/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertNote(note models.Note) error
}
