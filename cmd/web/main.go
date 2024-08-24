package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/devj1003/scribble/internal/config"
	"github.com/devj1003/scribble/internal/handlers"
	"github.com/devj1003/scribble/internal/render"
)

const portNumber = ":8080"

var App config.AppConfig

// Main is the main function
func main() {

	// App.UseCache = false

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	App.TemplateCache = tc

	repo := handlers.NewRepo(&App)
	handlers.NewHandlers(repo)

	http.HandleFunc("/", handlers.Repo.Home)

	fmt.Println(fmt.Printf("Staring application on port %s", portNumber))
	_ = http.ListenAndServe(portNumber, nil)
}
