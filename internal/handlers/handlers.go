package handlers

import (
	"net/http"

	"github.com/devj1003/scribble/internal/render"
)

// Home is the handler for the home page
func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "index.page.tmpl")
}
