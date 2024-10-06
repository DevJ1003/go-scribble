package repository

import "github.com/devj1003/scribble/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertNewUser(user models.User) error
	GetUserByID(id int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, testPassword string) (int, string, error)

	InsertNote(note models.Note) error
}
