package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/devj1003/scribble/internal/config"
	"github.com/devj1003/scribble/internal/drivers"
	"github.com/devj1003/scribble/internal/forms"
	"github.com/devj1003/scribble/internal/helpers"
	"github.com/devj1003/scribble/internal/models"
	"github.com/devj1003/scribble/internal/render"
	"github.com/devj1003/scribble/internal/repository"
	"github.com/devj1003/scribble/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo // for database
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *drivers.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewMysqlRepo(db.SQL, a), // for database
	}

}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Register is the handler for the register page
func (m *Repository) Register(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderShortTemplate(w, r, "register.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

// PostRegister is the handler for the registeration of the user
func (m *Repository) PostRegister(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		// log.Println(err)
		helpers.ServerError(w, err)
		return
	}

	createuser := models.User{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Password:  r.Form.Get("password"),
	}

	form := forms.New(r.PostForm)
	form.Required("first_name", "last_name", "email", "password")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["createuser"] = createuser

		// sending data to template
		render.RenderShortTemplate(w, r, "register.page.tmpl", &models.TemplateData{
			// Form: form,
			// Data: data,
		})

		return
	}

	// ==================================================================================
	// Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createuser.Password), 14)
	if err != nil {
		helpers.ServerError(w, err)
	}

	// sending hashedPassword to the user struct back
	createuser.Password = string(hashedPassword)
	// ==================================================================================

	err = m.DB.InsertNewUser(createuser)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "createuser", createuser)
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

// Login is the handler for the login page
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderShortTemplate(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})

}

// PostLogin handles logging the user in
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		render.RenderShortTemplate(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "success", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// Logout is the handler for logout
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {

	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	// ====================================================================
	userID := m.App.Session.Get(r.Context(), "user_id")
	if userID == nil {
		// Handle the case where user ID is not found in the session
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id, ok := userID.(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}
	// ====================================================================

	viewnoteatindex, err := m.DB.ViewNoteAtIndex(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["viewnoteatindex"] = viewnoteatindex

	// sending data to template
	render.RenderTemplate(w, r, "index.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

// CreateNewNote is the handler for the create-note page
func (m *Repository) CreateNewNote(w http.ResponseWriter, r *http.Request) {

	var emptyReservation models.Note
	data := make(map[string]interface{})
	data["createnote"] = emptyReservation

	// sending data to template
	render.RenderTemplate(w, r, "create-note.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})

}

// PostCreateNewNote is the handler for the create-note page
func (m *Repository) PostCreateNewNote(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		// log.Println(err)
		helpers.ServerError(w, err)
		return
	}

	createnote := models.Note{
		Title:   r.Form.Get("title"),
		Content: r.Form.Get("content"),
	}

	// ====================================================================
	userID := m.App.Session.Get(r.Context(), "user_id")
	if userID == nil {
		// Handle the case where user ID is not found in the session
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id, ok := userID.(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}
	// ====================================================================

	form := forms.New(r.PostForm)
	form.Required("title", "content")
	// form.Has("title", r)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["createnote"] = createnote

		// sending data to template
		render.RenderTemplate(w, r, "create-note.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	err = m.DB.InsertNote(createnote, id)
	if err != nil {
		helpers.ServerError(w, err)
	}

	m.App.Session.Put(r.Context(), "createnote", createnote)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// ViewNote is the handler for the create-note page
func (m *Repository) ViewNote(w http.ResponseWriter, r *http.Request) {

	nidStr := chi.URLParam(r, "nid")
	// Convert the string to an integer
	nid, err := strconv.Atoi(nidStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	// ====================================================================
	userID := m.App.Session.Get(r.Context(), "user_id")
	if userID == nil {
		// Handle the case where user ID is not found in the session
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id, ok := userID.(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}
	// ====================================================================

	viewnotedata, err := m.DB.ViewNoteData(nid, id)
	if err != nil {
		log.Println("Error fetching note data:", err)
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["viewnotedata"] = viewnotedata

	// sending data to template
	render.RenderTemplate(w, r, "view-note.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

// // EditNote is the handler for the edit-note page
// func (m *Repository) EditNote(w http.ResponseWriter, r *http.Request) {

// 	nidStr := chi.URLParam(r, "nid")
// 	// Convert the string to an integer
// 	nid, err := strconv.Atoi(nidStr)
// 	if err != nil {
// 		http.Error(w, "Invalid note ID", http.StatusBadRequest)
// 		return
// 	}

// 	// ====================================================================
// 	userID := m.App.Session.Get(r.Context(), "user_id")
// 	if userID == nil {
// 		// Handle the case where user ID is not found in the session
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}
// 	id, ok := userID.(int)
// 	if !ok {
// 		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
// 		return
// 	}
// 	// ====================================================================

// 	// exploded := strings.Split(r.RequestURI, "/")
// 	// src := exploded[3]
// 	// stringMap := make(map[string]string)
// 	// stringMap["src"] = src

// 	viewnotedata, err := m.DB.ViewNoteData(nid, id)
// 	if err != nil {
// 		log.Println("Error fetching note data:", err)
// 		helpers.ServerError(w, err)
// 		return
// 	}

// 	// viewnotedata.Title = r.Form.Get("title")
// 	// viewnotedata.Content = r.Form.Get("content")
// 	if r.Method == http.MethodPost {
// 		err := r.ParseForm()
// 		if err != nil {
// 			http.Error(w, "Unable to parse form", http.StatusBadRequest)
// 			return
// 		}
// 		// Now access form values
// 		viewnotedata.Title = r.Form.Get("title")
// 		viewnotedata.Content = r.Form.Get("content")
// 	}

// 	err = m.DB.UpdateNote(viewnotedata)
// 	if err != nil {
// 		helpers.ServerError(w, err)
// 		return
// 	}

// 	title := r.Form.Get("title")
// 	content := r.Form.Get("content")

// 	if content == "" {
// 		http.Redirect(w, r, fmt.Sprintf("/viewnote/%d", nid), http.StatusSeeOther)
// 	} else {
// 		http.Redirect(w, r, fmt.Sprintf("/viewnote?nid=%d&title=%s&content=%s", nid, title, content), http.StatusSeeOther)
// 	}

// 	// sending data to template
// 	// render.RenderTemplate(w, r, "edit-note.page.tmpl", &models.TemplateData{})
// 	// http.Redirect(w, r, "/", http.StatusSeeOther)

// }

func (m *Repository) EditNote(w http.ResponseWriter, r *http.Request) {
	// Fetching the note ID from the URL
	nidStr := chi.URLParam(r, "nid")
	nid, err := strconv.Atoi(nidStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	// Fetching the user ID from the session
	userID := m.App.Session.Get(r.Context(), "user_id")
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id, ok := userID.(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}

	// Fetch the note data
	viewnotedata, err := m.DB.ViewNoteData(nid, id)
	if err != nil {
		log.Println("Error fetching note data:", err)
		helpers.ServerError(w, err)
		return
	}

	// Handle form submission
	if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
			helpers.ServerError(w, err)
			return
		}

		// Update note data with form values
		viewnotedata.Title = r.Form.Get("title")
		viewnotedata.Content = r.Form.Get("content")

		// Now call UpdateNote with correct nid and id
		err = m.DB.UpdateNote(viewnotedata, nid, id)
		if err != nil {
			log.Println("Error updating note:", err)
			helpers.ServerError(w, err)
			return
		}

		// Redirect after update
		http.Redirect(w, r, fmt.Sprintf("/viewnote/%d", nid), http.StatusSeeOther)
		return
	}

	// Render the edit note page (GET request)
	data := make(map[string]interface{})
	data["editnotedata"] = viewnotedata

	editdata := &models.TemplateData{
		Data: data, // Pass the note data to the template
	}

	render.RenderTemplate(w, r, "edit-note.page.tmpl", editdata)
}

// DeleteNote is the handler to delete note
func (m *Repository) DeleteNote(w http.ResponseWriter, r *http.Request) {

	nidStr := chi.URLParam(r, "nid")
	// Convert the string to an integer
	nid, err := strconv.Atoi(nidStr)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}
	// ====================================================================
	userID := m.App.Session.Get(r.Context(), "user_id")
	if userID == nil {
		// Handle the case where user ID is not found in the session
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id, ok := userID.(int)
	if !ok {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}
	// ====================================================================
	err = m.DB.DeleteNote(nid, id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// m.App.Session.Put(r.Context(), "success", "Note deleted!")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// AdminDashboard is the handler for admin dashboard page
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})

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

// ProfilePage is the handler for the user-profile page
func (m *Repository) ProfilePage(w http.ResponseWriter, r *http.Request) {

	// sending data to template
	render.RenderTemplate(w, r, "user-profile-edit.page.tmpl", &models.TemplateData{})

}
