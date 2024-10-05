package models

import "time"

// User is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	Age         int
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Note is the note model
type Note struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
