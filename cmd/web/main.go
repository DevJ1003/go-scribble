package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/devj1003/scribble/internal/config"
	"github.com/devj1003/scribble/internal/drivers"
	"github.com/devj1003/scribble/internal/handlers"
	"github.com/devj1003/scribble/internal/helpers"
	"github.com/devj1003/scribble/internal/models"
	"github.com/devj1003/scribble/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main function
func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() (*drivers.DB, error) {

	// what am i going to put in the session
	gob.Register(models.User{})
	gob.Register(models.Note{})

	// changes this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to the database =========================================================
	fmt.Println("Connecting to the database...")
	db, err := drivers.ConnectSQL("root:password@tcp(127.0.0.1:3306)/scribble")
	if err != nil {
		log.Fatal("Cannot connect to the database, Dying...")
	}
	fmt.Println("Connected to database!")
	// defer db.SQL.Close()
	// =================================================================================

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	tc, err = render.CreateShortTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	// http.HandleFunc("/", handlers.Repo.Home)

	return db, nil

}
