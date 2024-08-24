package handlers

import (
	"net/http"

	"github.com/devj1003/scribble/internal/config"
	"github.com/devj1003/scribble/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository
var App *config.AppConfig

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "index.page.tmpl")
}
