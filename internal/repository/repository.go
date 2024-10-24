package repository

import "github.com/devj1003/scribble/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertNewUser(user models.User) error
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)

	InsertNote(note models.Note, id int) error
	ViewNoteAtIndex(id int) ([]models.Note, error)
	ViewNoteData(nid int, id int) (models.Note, error)
	UpdateNote(u models.Note, nid int, id int) error
	DeleteNote(nid int, id int) error
}
