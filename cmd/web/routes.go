package main

import (
	"net/http"

	"github.com/devj1003/scribble/internal/config"
	"github.com/devj1003/scribble/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// mux.Get("/", handlers.Repo.Home)

	mux.Get("/register", handlers.Repo.Register)
	mux.Post("/register", handlers.Repo.PostRegister)
	mux.Get("/login", handlers.Repo.Login)
	mux.Post("/login", handlers.Repo.PostLogin)
	mux.Get("/logout", handlers.Repo.Logout)

	mux.Get("/createnote", handlers.Repo.CreateNewNote)
	mux.Post("/createnote", handlers.Repo.PostCreateNewNote)
	mux.Get("/viewnote/{nid}", handlers.Repo.ViewNote)
	mux.Get("/editnote/{nid}", handlers.Repo.EditNote)
	mux.Get("/deletenote/{nid}", handlers.Repo.DeleteNote)
	mux.Get("/blank", handlers.Repo.Blank)
	mux.Get("/error500", handlers.Repo.Error500)
	mux.Get("/errorpage", handlers.Repo.Errorpage)
	mux.Get("/userlist", handlers.Repo.UserList)
	mux.Get("/profilepage", handlers.Repo.ProfilePage)

	fileServer := http.FileServer(http.Dir("../../static/"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	mux.Route("/", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/", handlers.Repo.Home)
	})

	return mux
}
