package handlers

import (
	"net/http"

	"github.com/devj1003/scribble/internal/config"
	"github.com/devj1003/scribble/internal/forms"
	"github.com/devj1003/scribble/internal/models"
	"github.com/devj1003/scribble/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

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
	// perform some logic

	// sending data to template
	render.RenderTemplate(w, r, "index.page.tmpl", &models.TemplateData{})

}

// Register is the handler for the register page
func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderShortTemplate(w, r, "register.page.tmpl", &models.TemplateData{})

}

// Login is the handler for the login page
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderTemplate(w, r, "login.page.tmpl", &models.TemplateData{})

}

// CreateNewNote is the handler for the create-note page
func (m *Repository) CreateNewNote(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderTemplate(w, r, "create-note.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

// PostCreateNewNote is the handler for the create-note page
func (m *Repository) PostCreateNewNote(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "create-note.page.tmpl", &models.TemplateData{})

}

// ViewNote is the handler for the create-note page
func (m *Repository) ViewNote(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderTemplate(w, r, "view-note.page.tmpl", &models.TemplateData{})

}

// EditNote is the handler for the edit-note page
func (m *Repository) EditNote(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderTemplate(w, r, "edit-note.page.tmpl", &models.TemplateData{})

}

// Blank is the handler for the blank page
func (m *Repository) Blank(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderTemplate(w, r, "pages-blank.page.tmpl", &models.TemplateData{})

}

// Error500 is the handler for the error500 page
func (m *Repository) Error500(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderTemplate(w, r, "pages-error-500.page.tmpl", &models.TemplateData{})

}

// Errorpage is the handler for the pages-error page
func (m *Repository) Errorpage(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderTemplate(w, r, "pages-error.page.tmpl", &models.TemplateData{})

}

// UserList is the handler for the user-list page
func (m *Repository) UserList(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderTemplate(w, r, "user-list.page.tmpl", &models.TemplateData{})

}

// ProfileEdit is the handler for the profile-edit page
func (m *Repository) ProfileEdit(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderTemplate(w, r, "user-profile-edit.page.tmpl", &models.TemplateData{})

}

// ProfilePage is the handler for the user-profile page
func (m *Repository) ProfilePage(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderTemplate(w, r, "user-profile.page.tmpl", &models.TemplateData{})

}
