package main

import (
	"fmt"
	"net/http"

	"github.com/devj1003/scribble/internal/handlers"
)

const portNumber = ":8080"

// Main is the main function
func main() {
	http.HandleFunc("/", handlers.Home)

	fmt.Println(fmt.Sprintf("Staring application on port %s", portNumber))
	_ = http.ListenAndServe(portNumber, nil)
}
